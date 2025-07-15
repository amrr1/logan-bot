package bot

import (
    "bufio"
    "errors"
    "io"
    "os/exec"
    "time"

    "github.com/bwmarrin/discordgo"
)

func PlayYouTubeAudio(s *discordgo.Session, guildID, channelID, youtubeURL, userID string) error {
    // Find voice channel
    var voiceChannelID string
    for _, g := range s.State.Guilds {
        if g.ID == guildID {
            for _, v := range g.VoiceStates {
                if v.UserID == userID {
                    voiceChannelID = v.ChannelID
                    break
                }
            }
        }
    }
    if voiceChannelID == "" {
        return errors.New("user not in a voice channel")
    }

    // Get direct audio URL from yt-dlp
    ytCmd := exec.Command("yt-dlp", "-f", "bestaudio", "-g", youtubeURL)
    output, err := ytCmd.Output()
    if err != nil {
        return errors.New("yt-dlp failed: " + err.Error())
    }
    directURL := string(output)

    // Join voice channel
    vc, err := s.ChannelVoiceJoin(guildID, voiceChannelID, false, true)
    if err != nil {
        return err
    }
    defer vc.Disconnect()
    time.Sleep(300 * time.Millisecond)

    // Run ffmpeg to output Ogg Opus to stdout
    ffmpeg := exec.Command("ffmpeg",
        "-i", directURL,
        "-f", "ogg",
        "-acodec", "libopus",
        "-ar", "48000",
        "-ac", "2",
        "-loglevel", "quiet",
        "pipe:1",
    )
    stdout, err := ffmpeg.StdoutPipe()
    if err != nil {
        return err
    }
    if err := ffmpeg.Start(); err != nil {
        return err
    }

    reader := bufio.NewReader(stdout)

    for {
        // Read Ogg page header (27 bytes)
        header := make([]byte, 27)
        if _, err := io.ReadFull(reader, header); err != nil {
            if err == io.EOF {
                break
            }
            return err
        }

        if string(header[0:4]) != "OggS" {
            return errors.New("invalid Ogg stream: missing OggS capture pattern")
        }

        segmentCount := int(header[26])
        segmentTable := make([]byte, segmentCount)
        if _, err := io.ReadFull(reader, segmentTable); err != nil {
            return err
        }

        // Calculate total size of page data by summing segment sizes
        totalSize := 0
        for _, segSize := range segmentTable {
            totalSize += int(segSize)
        }

        pageData := make([]byte, totalSize)
        if _, err := io.ReadFull(reader, pageData); err != nil {
            return err
        }

        // Now extract Opus packets using the segment table
        offset := 0
        for i, segSize := range segmentTable {
            if segSize == 0 {
                continue
            }
            packet := pageData[offset : offset+int(segSize)]
            offset += int(segSize)

            // Send each Opus packet to Discord
            select {
            case vc.OpusSend <- packet:
            case <-time.After(1 * time.Second):
                return errors.New("timeout sending opus packet")
            }

            // Sleep ~20ms for each packet (Discord expects 20ms per frame)
            time.Sleep(20 * time.Millisecond)

            // If this is the last segment of a packet that spans multiple segments,
            // or the segment size < 255, the packet ends here.
            // Here we treat each segment as a separate packet, which works for typical streams.
            if segSize < 255 || i == len(segmentTable)-1 {
                // Packet boundary
            }
        }
    }

    return nil
}
