
const vue = new Vue({
    el: '#polymerdesign-app',
    delimiters: ['{', '}'],
    data: {

    },
    mounted: () => {

    },
    methods: {
        toggleMenuButton: (e) => {
            let button = e.target

            if (document.activeElement === button) {
                setTimeout( () => {
                    button.blur()
                }, 0, button)
            }
        },
        hoverMenuButton: (e) => {
            let button = e.target
            let selected = document.activeElement

            if (document.activeElement.className && document.activeElement.className.includes('toolbar_menu_text')) {
                setTimeout( () => {
                    selected.blur()
                    button.focus()
                }, 0, button)
            }
        }
    }
})
//# sourceMappingURL=app.js.map