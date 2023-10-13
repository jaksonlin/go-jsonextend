package ast

type nodePlugins struct {
	plugins []ASTNodePlugin
}

func (plugin *nodePlugins) PreVisitPlugin(visitor JsonVisitor, node JsonNode) error {
	if len(plugin.plugins) == 0 {
		return nil
	}
	for _, plugin := range plugin.plugins {

		if err := plugin.PreVisitPlugin(visitor, node); err != nil {
			return err
		}
	}
	return nil
}

func (plugin *nodePlugins) PostVisitPlugin(visitor JsonVisitor, node JsonNode) error {
	if len(plugin.plugins) == 0 {
		return nil
	}
	for _, plugin := range plugin.plugins {
		if err := plugin.PostVisitPlugin(visitor, node); err != nil {
			return err
		}
	}
	return nil
}

func (plugin *nodePlugins) RegisterPlugin(p ASTNodePlugin) {
	if plugin.plugins == nil {
		plugin.plugins = make([]ASTNodePlugin, 0)
	}
	for _, item := range plugin.plugins {
		if item.PluginName() == p.PluginName() {
			return
		}
	}
	plugin.plugins = append(plugin.plugins, p)
}

func (plugin *nodePlugins) RegisterPluginAtFirst(p ASTNodePlugin) {
	if plugin.plugins == nil {
		plugin.plugins = make([]ASTNodePlugin, 0)
		plugin.plugins = append(plugin.plugins, p)
		return
	}
	for i, item := range plugin.plugins {
		if item.PluginName() == p.PluginName() {
			plugin.plugins = append(plugin.plugins[:i], plugin.plugins[i+1:]...)
			break
		}
	}
	plugin.plugins = append([]ASTNodePlugin{p}, plugin.plugins...)
}

func (plugin *nodePlugins) RemovePlugin(name string) {
	for i, p := range plugin.plugins {
		if p.PluginName() == name {
			plugin.plugins = append(plugin.plugins[:i], plugin.plugins[i+1:]...)
			return
		}
	}
}
