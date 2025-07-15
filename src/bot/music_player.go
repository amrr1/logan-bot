package bot

import (
    "errors"
    "os/exec"
    "time"

    "github.com/bwmarrin/discordgo"
)

// PlayYouTubeAudio joins the user's voice channel and plays audio from a YouTube link
func PlayYouTubeAudio(s *discordgo.Session, guildID, channelID, youtubeURL string, userID string) error {
    var voiceChannelID string
    for _, vs := range s.State.Guilds {
        if vs.ID == guildID {
            for _, v := range vs.VoiceStates {
                if v.UserID == userID { // Use the command user's ID
                    voiceChannelID = v.ChannelID
                    break
                }
            }
        }
    }
    if voiceChannelID == "" {
        return errors.New("user not in a voice channel")
    }

    vc, err := s.ChannelVoiceJoin(guildID, voiceChannelID, false, true)
    if err != nil {
        return err
    }
    defer vc.Disconnect()

    // Download and convert YouTube audio to PCM using yt-dlp and ffmpeg
    // Requires yt-dlp and ffmpeg installed on the system
    cmd := exec.Command("yt-dlp", "-f", "bestaudio", "--extract-audio", "--audio-format", "wav", "-o", "audio.wav", youtubeURL)
    if err := cmd.Run(); err != nil {
        return err
    }

    // Play the audio file (audio.wav)
    audioFile := "audio.wav"
    err = playAudioFile(vc, audioFile)
    if err != nil {
        return err
    }

    // Clean up
    exec.Command("rm", audioFile).Run()
    return nil
}

// playAudioFile streams a wav file to Discord (simplified, not production ready)
func playAudioFile(vc *discordgo.VoiceConnection, filename string) error {
    // This is a stub. Implement PCM streaming here.
    // For a real implementation, use github.com/bwmarrin/dca or similar.
    time.Sleep(10 * time.Second) // Simulate playback
    return nil
}