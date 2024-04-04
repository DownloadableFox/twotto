package core

import "github.com/bwmarrin/discordgo"

type CommandStack struct {
	commands []*discordgo.ApplicationCommand
}

func NewCommandStack() *CommandStack {
	return &CommandStack{
		commands: make([]*discordgo.ApplicationCommand, 0),
	}
}

func (c *CommandStack) AddCommand(command *discordgo.ApplicationCommand) {
	c.commands = append(c.commands, command)
}

func (c *CommandStack) For(session *discordgo.Session, guildId string) error {
	for _, command := range c.commands {
		if _, err := session.ApplicationCommandCreate(session.State.User.ID, guildId, command); err != nil {
			return err
		}
	}

	return nil
}

func ApplyCommands(commands ...*discordgo.ApplicationCommand) *CommandStack {
	stack := NewCommandStack()
	for _, command := range commands {
		stack.AddCommand(command)
	}
	return stack
}
