package ast

type nodePlugins struct {
	preVisitPlugins  []ASTNodePlugin
	postVisitPlugins []ASTNodePlugin
}

func (plugin *nodePlugins) PreVisitPlugin(visitor JsonVisitor, node JsonNode) error {
	if node.IsVisited() {
		return nil
	}

	for _, p := range plugin.preVisitPlugins {
		if err := p.PreVisitPlugin(visitor, node); err != nil {
			return err
		}
	}
	return nil
}

func (plugin *nodePlugins) PostVisitPlugin(visitor JsonVisitor, node JsonNode) error {

	for _, p := range plugin.postVisitPlugins {
		if err := p.PostVisitPlugin(visitor, node); err != nil {
			return err
		}
	}
	return nil
}

func (plugin *nodePlugins) RegisterPrevisitPlugin(p ASTNodePlugin) {
	if plugin.preVisitPlugins == nil {
		plugin.preVisitPlugins = make([]ASTNodePlugin, 0)
	}
	for _, item := range plugin.preVisitPlugins {
		if item.PluginName() == p.PluginName() {
			return
		}
	}
	plugin.preVisitPlugins = append(plugin.preVisitPlugins, p)
}

func (plugin *nodePlugins) RegisterPostvisitPlugin(p ASTNodePlugin) {
	if plugin.postVisitPlugins == nil {
		plugin.postVisitPlugins = make([]ASTNodePlugin, 0)
	}
	for _, item := range plugin.postVisitPlugins {
		if item.PluginName() == p.PluginName() {
			return
		}
	}
	plugin.postVisitPlugins = append(plugin.postVisitPlugins, p)
}

func (plugin *nodePlugins) GetPrevisitPlugins() []ASTNodePlugin {
	return plugin.preVisitPlugins
}

func (plugin *nodePlugins) GetPostvisitPlugins() []ASTNodePlugin {
	return plugin.postVisitPlugins
}

func (plugin *nodePlugins) RemovePrevisitPlugin(name string) {
	for i, p := range plugin.preVisitPlugins {
		if p.PluginName() == name {
			plugin.preVisitPlugins = append(plugin.preVisitPlugins[:i], plugin.preVisitPlugins[i+1:]...)
			return
		}
	}
}

func (plugin *nodePlugins) RemovePostvisitPlugin(name string) {
	for i, p := range plugin.postVisitPlugins {
		if p.PluginName() == name {
			plugin.postVisitPlugins = append(plugin.postVisitPlugins[:i], plugin.postVisitPlugins[i+1:]...)
			return
		}
	}
}
