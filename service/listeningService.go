package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
)

func (s *CommandService) ListeningService(response http.ResponseWriter, request *http.Request) {
	options := s.discordMessage.Data.Options[0]
	msg := ""
	requiresUpdate := false

	if options.Value.(bool) && strings.Contains(s.discordMessage.Member.Nick, utils.NICKNAME_SUFFIX) {
		msg = "You are already set to listen."
	} else if !options.Value.(bool) && !strings.Contains(s.discordMessage.Member.Nick, utils.NICKNAME_SUFFIX) {
		msg = "Your nickname remains unchanged."
	} else {
		requiresUpdate = true
		msg = "Your nickname will be updated shortly."
	}

	if requiresUpdate {
		dataPacket := dtos.DataPacket{
			UserID:      s.discordMessage.Member.User.ID,
			CommandName: utils.CommandNames.Listening,
			MetaData: map[string]string{
				"value":    fmt.Sprint(options.Value),
				"nickname": s.discordMessage.Member.Nick,
			},
		}
		bytePacket, err := dtos.ToByte(&dataPacket)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := queue.SendMessage(bytePacket); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	messageResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	utils.Success.NewDiscordResponse(response, "Success", messageResponse)
}
