package bot

import (
    "strings"

    "github.com/bwmarrin/discordgo"
)

// MessageHandler handles incoming Discord messages
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }

    content := strings.ToLower(m.Content)

    if strings.Contains(content, "logan") {
        s.ChannelMessageSend(m.ChannelID, "IM LOGING OUT")
    }

    if strings.HasPrefix(content, "!play ") {
        parts := strings.Fields(m.Content)
        if len(parts) >= 2 {
            youtubeURL := parts[1]
            go func() {
				err := PlayYouTubeAudio(s, m.GuildID, m.ChannelID, youtubeURL, m.Author.ID)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Error playing audio: "+err.Error())
				}
            }()
        } else {
            s.ChannelMessageSend(m.ChannelID, "Usage: !play <youtube_link>")
        }
    }
}