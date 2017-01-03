import sublime_plugin


class FoldCommentsCommand(sublime_plugin.TextCommand):

    def run(self, edit):
        self.view.fold(self.view.find_by_selector('comment'))


class UnfoldCommentsCommand(sublime_plugin.TextCommand):

    def run(self, edit):
        for region in self.view.find_by_selector('comment'):
            self.view.unfold(region)
