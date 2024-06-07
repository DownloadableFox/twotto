package core

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type CommandBuilder struct {
	discordgo.ApplicationCommand
}

func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{}
}

func (c *CommandBuilder) SetName(name string) *CommandBuilder {
	c.Name = name
	return c
}

func (c *CommandBuilder) SetDescription(description string) *CommandBuilder {
	c.Description = description
	return c
}

func (c *CommandBuilder) SetDMPermission(permission bool) *CommandBuilder {
	c.DMPermission = &permission
	return c
}

func (c *CommandBuilder) SetDefaultMemberPermissions(permissions int64) *CommandBuilder {
	c.DefaultMemberPermissions = &permissions
	return c
}

func (c *CommandBuilder) SetNSFW(nsfw bool) *CommandBuilder {
	c.NSFW = &nsfw
	return c
}

func (c *CommandBuilder) AddNameLocale(locale discordgo.Locale, name string) *CommandBuilder {
	if c.NameLocalizations == nil {
		temp := make(map[discordgo.Locale]string)
		c.NameLocalizations = &temp
	}

	(*c.NameLocalizations)[locale] = name

	return c
}

func (c *CommandBuilder) AddDescriptionLocale(locale discordgo.Locale, description string) *CommandBuilder {
	if c.DescriptionLocalizations == nil {
		temp := make(map[discordgo.Locale]string)
		c.DescriptionLocalizations = &temp
	}

	(*c.DescriptionLocalizations)[locale] = description

	return c
}

type StringOptionSupplier func(*StringOptionBuilder)
type IntegerOptionSupplier func(*IntegerOptionBuilder)
type BooleanOptionSupplier func(*BooleanOptionBuilder)
type UserOptionSupplier func(*UserOptionBuilder)
type ChannelOptionSupplier func(*ChannelOptionBuilder)
type RoleOptionSupplier func(*RoleOptionBuilder)
type MentionableOptionSupplier func(*MentionableOptionBuilder)
type NumberOptionSupplier func(*NumberOptionBuilder)

type CommandOption interface {
	Build() *discordgo.ApplicationCommandOption
}

func (c *CommandBuilder) AddOption(option CommandOption) *CommandBuilder {
	if c.Options == nil {
		temp := make([]*discordgo.ApplicationCommandOption, 0)
		c.Options = temp
	}

	c.Options = append(c.Options, option.Build())
	return c
}

func (c *CommandBuilder) AddStringOption(supplier StringOptionSupplier) *CommandBuilder {
	builder := NewStringOptionBuilder()
	supplier(builder)
	c.AddOption(builder)
	return c
}

func (c *CommandBuilder) AddIntegerOption(supplier IntegerOptionSupplier) *CommandBuilder {
	builder := NewIntegerOptionBuilder()
	supplier(builder)
	c.AddOption(builder)
	return c
}

func (c *CommandBuilder) AddBooleanOption(supplier BooleanOptionSupplier) *CommandBuilder {
	builder := NewBooleanOptionBuilder()
	supplier(builder)
	c.AddOption(builder)
	return c
}

func (c *CommandBuilder) AddUserOption(supplier UserOptionSupplier) *CommandBuilder {
	builder := NewUserOptionBuilder()
	supplier(builder)
	c.AddOption(builder)
	return c
}

func (c *CommandBuilder) AddChannelOption(supplier ChannelOptionSupplier) *CommandBuilder {
	builder := NewChannelOptionBuilder()
	supplier(builder)
	c.AddOption(builder)
	return c
}

func (c *CommandBuilder) AddRoleOption(supplier RoleOptionSupplier) *CommandBuilder {
	builder := NewRoleOptionBuilder()
	supplier(builder)
	c.AddOption(builder)
	return c
}

func (c *CommandBuilder) AddMentionableOption(supplier MentionableOptionSupplier) *CommandBuilder {
	builder := NewMentionableOptionBuilder()
	supplier(builder)
	c.AddOption(builder)
	return c
}

func (c *CommandBuilder) AddNumberOption(supplier NumberOptionSupplier) *CommandBuilder {
	builder := NewNumberOptionBuilder()
	supplier(builder)
	c.AddOption(builder)
	return c
}

func (c *CommandBuilder) AddSubCommand(supplier func(*SubCommandBuilder)) *CommandBuilder {
	builder := NewSubCommandBuilder()
	supplier(builder)
	c.AddOption(builder)
	return c
}

func (c *CommandBuilder) Build() *discordgo.ApplicationCommand {
	return &c.ApplicationCommand
}

type StringOptionBuilder struct {
	discordgo.ApplicationCommandOption
}

func NewStringOptionBuilder() *StringOptionBuilder {
	return &StringOptionBuilder{
		ApplicationCommandOption: discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionString,
		},
	}
}

func (s *StringOptionBuilder) SetName(name string) *StringOptionBuilder {
	s.Name = name
	return s
}

func (s *StringOptionBuilder) SetDescription(description string) *StringOptionBuilder {
	s.Description = description
	return s
}

func (s *StringOptionBuilder) SetRequired(required bool) *StringOptionBuilder {
	s.Required = required
	return s
}

func (s *StringOptionBuilder) Build() *discordgo.ApplicationCommandOption {
	return &s.ApplicationCommandOption
}

type IntegerOptionBuilder struct {
	discordgo.ApplicationCommandOption
}

func NewIntegerOptionBuilder() *IntegerOptionBuilder {
	return &IntegerOptionBuilder{
		ApplicationCommandOption: discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionInteger,
		},
	}
}

func (i *IntegerOptionBuilder) SetName(name string) *IntegerOptionBuilder {
	i.Name = name
	return i
}

func (i *IntegerOptionBuilder) SetDescription(description string) *IntegerOptionBuilder {
	i.Description = description
	return i
}

func (i *IntegerOptionBuilder) SetRequired(required bool) *IntegerOptionBuilder {
	i.Required = required
	return i
}

func (i *IntegerOptionBuilder) Build() *discordgo.ApplicationCommandOption {
	return &i.ApplicationCommandOption
}

type BooleanOptionBuilder struct {
	discordgo.ApplicationCommandOption
}

func NewBooleanOptionBuilder() *BooleanOptionBuilder {
	return &BooleanOptionBuilder{
		ApplicationCommandOption: discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionBoolean,
		},
	}
}

func (b *BooleanOptionBuilder) SetName(name string) *BooleanOptionBuilder {
	b.Name = name
	return b
}

func (b *BooleanOptionBuilder) SetDescription(description string) *BooleanOptionBuilder {
	b.Description = description
	return b
}

func (b *BooleanOptionBuilder) SetRequired(required bool) *BooleanOptionBuilder {
	b.Required = required
	return b
}

func (b *BooleanOptionBuilder) Build() *discordgo.ApplicationCommandOption {
	return &b.ApplicationCommandOption
}

type UserOptionBuilder struct {
	discordgo.ApplicationCommandOption
}

func NewUserOptionBuilder() *UserOptionBuilder {
	return &UserOptionBuilder{
		ApplicationCommandOption: discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionUser,
		},
	}
}

func (u *UserOptionBuilder) SetName(name string) *UserOptionBuilder {
	u.Name = name
	return u
}

func (u *UserOptionBuilder) SetDescription(description string) *UserOptionBuilder {
	u.Description = description
	return u
}

func (u *UserOptionBuilder) SetRequired(required bool) *UserOptionBuilder {
	u.Required = required
	return u
}

func (u *UserOptionBuilder) Build() *discordgo.ApplicationCommandOption {
	return &u.ApplicationCommandOption
}

type ChannelOptionBuilder struct {
	discordgo.ApplicationCommandOption
}

func NewChannelOptionBuilder() *ChannelOptionBuilder {
	return &ChannelOptionBuilder{
		ApplicationCommandOption: discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionChannel,
		},
	}
}

func (c *ChannelOptionBuilder) SetName(name string) *ChannelOptionBuilder {
	c.Name = name
	return c
}

func (c *ChannelOptionBuilder) SetDescription(description string) *ChannelOptionBuilder {
	c.Description = description
	return c
}

func (c *ChannelOptionBuilder) SetRequired(required bool) *ChannelOptionBuilder {
	c.Required = required
	return c
}

func (c *ChannelOptionBuilder) Build() *discordgo.ApplicationCommandOption {
	return &c.ApplicationCommandOption
}

type RoleOptionBuilder struct {
	discordgo.ApplicationCommandOption
}

func NewRoleOptionBuilder() *RoleOptionBuilder {
	return &RoleOptionBuilder{
		ApplicationCommandOption: discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionRole,
		},
	}
}

func (r *RoleOptionBuilder) SetName(name string) *RoleOptionBuilder {
	r.Name = name
	return r
}

func (r *RoleOptionBuilder) SetDescription(description string) *RoleOptionBuilder {
	r.Description = description
	return r
}

func (r *RoleOptionBuilder) SetRequired(required bool) *RoleOptionBuilder {
	r.Required = required
	return r
}

func (r *RoleOptionBuilder) Build() *discordgo.ApplicationCommandOption {
	return &r.ApplicationCommandOption
}

type MentionableOptionBuilder struct {
	discordgo.ApplicationCommandOption
}

func NewMentionableOptionBuilder() *MentionableOptionBuilder {
	return &MentionableOptionBuilder{
		ApplicationCommandOption: discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionMentionable,
		},
	}
}

func (m *MentionableOptionBuilder) SetName(name string) *MentionableOptionBuilder {
	m.Name = name
	return m
}

func (m *MentionableOptionBuilder) SetDescription(description string) *MentionableOptionBuilder {
	m.Description = description
	return m
}

func (m *MentionableOptionBuilder) SetRequired(required bool) *MentionableOptionBuilder {
	m.Required = required
	return m
}

func (m *MentionableOptionBuilder) Build() *discordgo.ApplicationCommandOption {
	return &m.ApplicationCommandOption
}

type NumberOptionBuilder struct {
	discordgo.ApplicationCommandOption
}

func NewNumberOptionBuilder() *NumberOptionBuilder {
	return &NumberOptionBuilder{
		ApplicationCommandOption: discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionNumber,
		},
	}
}

func (n *NumberOptionBuilder) SetName(name string) *NumberOptionBuilder {
	n.Name = name
	return n
}

func (n *NumberOptionBuilder) SetDescription(description string) *NumberOptionBuilder {
	n.Description = description
	return n
}

func (n *NumberOptionBuilder) SetRequired(required bool) *NumberOptionBuilder {
	n.Required = required
	return n
}

func (n *NumberOptionBuilder) SetMin(min float64) *NumberOptionBuilder {
	n.MinValue = &min
	return n
}

func (n *NumberOptionBuilder) SetMax(max float64) *NumberOptionBuilder {
	n.MaxValue = max
	return n
}

func (n *NumberOptionBuilder) Build() *discordgo.ApplicationCommandOption {
	return &n.ApplicationCommandOption
}

type OptionBuilder struct {
	discordgo.ApplicationCommandOption
}

type SubCommandBuilder struct {
	discordgo.ApplicationCommandOption
}

func NewSubCommandBuilder() *SubCommandBuilder {
	return &SubCommandBuilder{
		ApplicationCommandOption: discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionSubCommand,
		},
	}
}

func (s *SubCommandBuilder) SetName(name string) *SubCommandBuilder {
	s.Name = name
	return s
}

func (s *SubCommandBuilder) SetDescription(description string) *SubCommandBuilder {
	s.Description = description
	return s
}

func (s *SubCommandBuilder) AddOption(option CommandOption) *SubCommandBuilder {
	if s.Options == nil {
		temp := make([]*discordgo.ApplicationCommandOption, 0)
		s.Options = temp
	}

	s.Options = append(s.Options, option.Build())
	return s
}

func (s *SubCommandBuilder) AddStringOption(supplier StringOptionSupplier) *SubCommandBuilder {
	builder := NewStringOptionBuilder()
	supplier(builder)
	s.AddOption(builder)
	return s
}

func (s *SubCommandBuilder) AddIntegerOption(supplier IntegerOptionSupplier) *SubCommandBuilder {
	builder := NewIntegerOptionBuilder()
	supplier(builder)
	s.AddOption(builder)
	return s
}

func (s *SubCommandBuilder) AddBooleanOption(supplier BooleanOptionSupplier) *SubCommandBuilder {
	builder := NewBooleanOptionBuilder()
	supplier(builder)
	s.AddOption(builder)
	return s
}

func (s *SubCommandBuilder) AddUserOption(supplier UserOptionSupplier) *SubCommandBuilder {
	builder := NewUserOptionBuilder()
	supplier(builder)
	s.AddOption(builder)
	return s
}

func (s *SubCommandBuilder) AddChannelOption(supplier ChannelOptionSupplier) *SubCommandBuilder {
	builder := NewChannelOptionBuilder()
	supplier(builder)
	s.AddOption(builder)
	return s
}

func (s *SubCommandBuilder) AddRoleOption(supplier RoleOptionSupplier) *SubCommandBuilder {
	builder := NewRoleOptionBuilder()
	supplier(builder)
	s.AddOption(builder)
	return s
}

func (s *SubCommandBuilder) AddMentionableOption(supplier MentionableOptionSupplier) *SubCommandBuilder {
	builder := NewMentionableOptionBuilder()
	supplier(builder)
	s.AddOption(builder)
	return s
}

func (s *SubCommandBuilder) AddNumberOption(supplier NumberOptionSupplier) *SubCommandBuilder {
	builder := NewNumberOptionBuilder()
	supplier(builder)
	s.AddOption(builder)
	return s
}

func (s *SubCommandBuilder) Build() *discordgo.ApplicationCommandOption {
	return &s.ApplicationCommandOption
}

// Command data helper
var ErrOptionNotFound = errors.New("option not found")
var ErrOptionUnexpectedType = errors.New("unexpected option type")

func GetStringOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string) (string, error) {
	for _, option := range data {
		if option.Name == name {
			if option.Type != discordgo.ApplicationCommandOptionString {
				return "", ErrOptionUnexpectedType
			}

			return option.StringValue(), nil
		}
	}

	return "", ErrOptionNotFound
}

func GetStringDefaultOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string, def string) string {
	value, err := GetStringOption(data, name)
	if err != nil {
		return def
	}

	return value
}

func GetIntegerOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string) (int, error) {
	for _, option := range data {
		if option.Name == name {
			if option.Type != discordgo.ApplicationCommandOptionInteger {
				return 0, ErrOptionUnexpectedType
			}

			return int(option.IntValue()), nil
		}
	}

	return 0, ErrOptionNotFound
}

func GetIntegerDefaultOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string, def int) int {
	value, err := GetIntegerOption(data, name)
	if err != nil {
		return def
	}

	return value
}

func GetBooleanOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string) (bool, error) {
	for _, option := range data {
		if option.Name == name {
			if option.Type != discordgo.ApplicationCommandOptionBoolean {
				return false, ErrOptionUnexpectedType
			}

			return option.BoolValue(), nil
		}
	}

	return false, ErrOptionNotFound
}

func GetBooleanDefaultOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string, def bool) bool {
	value, err := GetBooleanOption(data, name)
	if err != nil {
		return def
	}

	return value
}

/*
func GetUserOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string) (*discordgo.User, error) {
	for _, option := range data {
		if option.Name == name {
			if option.Type != discordgo.ApplicationCommandOptionUser {
				return nil, ErrOptionUnexpectedType
			}

			return option.UserValue(), nil
		}
	}

	return nil, ErrOptionNotFound
}

func GetUserDefaultOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string, def *discordgo.User) *discordgo.User {
	value, err := GetUserOption(data, name)
	if err != nil {
		return def
	}

	return value
}

func GetChannelOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string) (*discordgo.Channel, error) {
	for _, option := range data {
		if option.Name == name {
			if option.Type != discordgo.ApplicationCommandOptionChannel {
				return nil, ErrOptionUnexpectedType
			}

			return option.ChannelValue(), nil
		}
	}

	return nil, ErrOptionNotFound
}

func GetChannelDefaultOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string, def *discordgo.Channel) *discordgo.Channel {
	value, err := GetChannelOption(data, name)
	if err != nil {
		return def
	}

	return value
}

func GetRoleOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string, session *discordgo.Session, guild string) (*discordgo.Role, error) {
	for _, option := range data {
		if option.Name == name {
			if option.Type != discordgo.ApplicationCommandOptionRole {
				return nil, ErrOptionUnexpectedType
			}

			return option.RoleValue(session, guild), nil
		}
	}

	return nil, ErrOptionNotFound
}

func GetRoleDefaultOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string, def *discordgo.Role) *discordgo.Role {
	value, err := GetRoleOption(data, )
	if err != nil {
		return def
	}

	return value
}
*/

func GetMentionableOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string) (string, error) {
	for _, option := range data {
		if option.Name == name {
			if option.Type != discordgo.ApplicationCommandOptionMentionable {
				return "", ErrOptionUnexpectedType
			}

			return option.StringValue(), nil
		}
	}

	return "", ErrOptionNotFound
}

func GetMentionableDefaultOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string, def string) string {
	value, err := GetMentionableOption(data, name)
	if err != nil {
		return def
	}

	return value
}

func GetNumberOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string) (float64, error) {
	for _, option := range data {
		if option.Name == name {
			if option.Type != discordgo.ApplicationCommandOptionNumber {
				return 0, ErrOptionUnexpectedType
			}

			return option.FloatValue(), nil
		}
	}

	return 0, ErrOptionNotFound
}

func GetNumberDefaultOption(data []*discordgo.ApplicationCommandInteractionDataOption, name string, def float64) float64 {
	value, err := GetNumberOption(data, name)
	if err != nil {
		return def
	}

	return value
}

func GetSubcommandOption(data []*discordgo.ApplicationCommandInteractionDataOption) string {
	if len(data) == 0 {
		return ""
	}

	return data[0].Name
}
