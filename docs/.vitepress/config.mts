import {defineConfig} from 'vitepress'

export default defineConfig({
	title: 'posh',
	description: 'An interactive, isolated and hackable Makefile',
	lang: 'en-US',
	lastUpdated: true,
	appearance: 'dark',
	ignoreDeadLinks: true,
	base: '/posh/',
	sitemap: {
		hostname: 'https://foomo.github.io/posh',
	},
	themeConfig: {
		logo: '/logo.png',
		outline: [2, 4],
		nav: [
			{text: 'Guide', link: '/guide/introduction'},
			{text: 'Usage', link: '/usage/prompt'},
			{text: 'CLI', link: '/reference/cli/'},
			{text: 'Plugin', link: '/plugin/overview'},
			{text: 'Recipes', link: '/recipes/local-dev-servers'},
		],
		sidebar: [
			{
				text: 'Guide',
				items: [
					{text: 'Introduction', link: '/guide/introduction'},
					{text: 'Installation', link: '/guide/installation'},
					{text: 'Quick Start', link: '/guide/quick-start'},
					{text: 'Concepts', link: '/guide/concepts'},
				],
			},
			{
				text: 'Usage',
				collapsed: true,
				items: [
					{text: 'The Interactive Prompt', link: '/usage/prompt'},
					{text: 'Built-in Commands', link: '/usage/builtins'},
					{text: 'Configuration', link: '/usage/configuration'},
					{text: 'Project Layout', link: '/usage/layout'},
				],
			},
			{
				text: 'CLI Reference',
				collapsed: true,
				items: [
					{text: 'Overview', link: '/reference/cli/'},
					{text: 'posh', link: '/reference/cli/posh'},
					{text: 'posh init', link: '/reference/cli/posh_init'},
					{text: 'posh config', link: '/reference/cli/posh_config'},
					{text: 'posh version', link: '/reference/cli/posh_version'},
					{text: 'posh prompt', link: '/reference/cli/posh_prompt'},
					{text: 'posh execute', link: '/reference/cli/posh_execute'},
					{text: 'posh require', link: '/reference/cli/posh_require'},
					{text: 'posh brew', link: '/reference/cli/posh_brew'},
				],
			},
			{
				text: 'Plugin Authoring',
				collapsed: true,
				items: [
					{text: 'Overview', link: '/plugin/overview'},
					{text: 'Writing Commands', link: '/plugin/writing-commands'},
					{text: 'Integrations', link: '/plugin/integrations'},
				],
			},
			{
				text: 'Recipes',
				collapsed: true,
				items: [
					{text: 'Local Dev Servers', link: '/recipes/local-dev-servers'},
				],
			},
			{
				text: 'Contributing',
				collapsed: true,
				items: [
					{text: 'Guideline', link: '/CONTRIBUTING'},
					{text: 'Code of conduct', link: '/CODE_OF_CONDUCT'},
					{text: 'Security guidelines', link: '/SECURITY'},
				],
			}
		],
		socialLinks: [
			{icon: 'github', link: 'https://github.com/foomo/posh'},
		],
		editLink: {
			pattern: 'https://github.com/foomo/posh/edit/main/docs/:path',
		},
		search: {
			provider: 'local',
		},
		footer: {
			message: 'Made with ♥ <a href="https://www.foomo.org">foomo</a> by <a href="https://www.bestbytes.com">bestbytes</a>',
		},
	},
	markdown: {
		theme: {
			light: 'catppuccin-latte',
			dark: 'catppuccin-frappe',
		},
	},
	head: [
		['meta', {name: 'theme-color', content: '#ffffff'}],
		['link', {rel: 'icon', href: '/logo.png'}],
		['meta', {name: 'author', content: 'foomo by bestbytes'}],
		['meta', {property: 'og:title', content: 'foomo/posh'}],
		[
			'meta',
			{
				property: 'og:image',
				content: 'https://github.com/foomo/posh/blob/main/docs/public/banner.png?raw=true',
			},
		],
		[
			'meta',
			{
				property: 'og:description',
				content: 'An interactive, isolated and hackable Makefile.',
			},
		],
		['meta', {name: 'twitter:card', content: 'summary_large_image'}],
		[
			'meta',
			{
				name: 'twitter:image',
				content: 'https://github.com/foomo/posh/blob/main/docs/public/banner.png?raw=true',
			},
		],
		[
			'meta',
			{
				name: 'viewport',
				content: 'width=device-width, initial-scale=1.0, viewport-fit=cover',
			},
		],
	],
})
