package debug

import "github.com/downloadablefox/twotto/core"

type ModuleProvider interface {
	// Modules
	EnableModule(module string, guildId string) error
	DisableModule(module string, guildId string) error
	IsModuleEnabled(module string, guildId string) bool

	// Features
	EnableFeature(identifier *core.Identifier, guildId string) error
	DisableFeature(identifier *core.Identifier, guildId string) error
	IsFeatureEnabled(identifier *core.Identifier, guildId string) bool
}

type GuildData struct {
	enabledModules  map[string]bool
	enabledFeatures map[string]bool
}

type InMemoryModuleProvider struct {
	guildData map[string]*GuildData
}

func NewInMemoryModuleProvider() ModuleProvider {
	return &InMemoryModuleProvider{
		guildData: make(map[string]*GuildData),
	}
}

func (p *InMemoryModuleProvider) GetGuildData(guildId string) *GuildData {
	data, ok := p.guildData[guildId]
	if !ok {
		p.guildData[guildId] = &GuildData{
			enabledModules:  make(map[string]bool),
			enabledFeatures: make(map[string]bool),
		}

		return p.guildData[guildId]
	}

	return data
}
func (p *InMemoryModuleProvider) EnableModule(module string, guildId string) error {
	data := p.GetGuildData(guildId)
	data.enabledModules[module] = true

	return nil
}

func (p *InMemoryModuleProvider) DisableModule(module string, guildId string) error {
	data := p.GetGuildData(guildId)
	data.enabledModules[module] = false

	return nil
}

func (p *InMemoryModuleProvider) IsModuleEnabled(module string, guildId string) bool {
	data := p.GetGuildData(guildId)
	return data.enabledModules[module]
}

func (p *InMemoryModuleProvider) EnableFeature(identifier *core.Identifier, guildId string) error {
	data := p.GetGuildData(guildId)
	data.enabledFeatures[identifier.String()] = true

	return nil
}

func (p *InMemoryModuleProvider) DisableFeature(identifier *core.Identifier, guildId string) error {
	data := p.GetGuildData(guildId)
	data.enabledFeatures[identifier.String()] = false

	return nil
}

func (p *InMemoryModuleProvider) IsFeatureEnabled(identifier *core.Identifier, guildId string) bool {
	data := p.GetGuildData(guildId)
	return data.enabledFeatures[identifier.String()]
}
