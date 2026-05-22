export type Theme = 'dark' | 'light' | 'system';
let currentTheme = $state<Theme>('system');
export const themeState = {
	get current() { return currentTheme; },
	set current(val: Theme) {
		currentTheme = val;
		localStorage.setItem('wat-theme', val);
		applyTheme(val);
	}
};
export function initTheme() {
	const saved = localStorage.getItem('wat-theme') as Theme | null;
	const theme = saved || 'system';
	currentTheme = theme;
	applyTheme(theme);
	const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
	mediaQuery.addEventListener('change', () => {
		if (currentTheme === 'system') {
			applyTheme('system');
		}
	});
}
function applyTheme(theme: Theme) {
	const root = document.documentElement;
	let isDark = false;
	if (theme === 'system') {
		isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
	} else {
		isDark = theme === 'dark';
	}
	if (isDark) {
		root.classList.add('dark');
	} else {
		root.classList.remove('dark');
	}
}
