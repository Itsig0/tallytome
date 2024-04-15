// HTMX event listener
document.addEventListener("DOMContentLoaded", (event) => {
    document.body.addEventListener("htmx:beforeSwap", function(evt) {
        if (evt.detail.xhr.status === 422) {
            evt.detail.shouldSwap = true;
            evt.detail.isError = false;
        }
    });

    document.body.addEventListener("HPUpdated", function() {
        var input = document.getElementById("damageInput");
        input.value = "";
        input.removeAttribute("aria-invalid");

        var errorMessage = document.getElementById("damageError");
        if (errorMessage) {
            errorMessage.parentNode.removeChild(errorMessage);
        }

        var savingthrow = document.getElementById("savingthrow");
        savingthrow.checked = false;
    });

    document.body.addEventListener("ManaUpdated", function() {
        var input = document.getElementById("manaInput");
        input.value = "";
        input.removeAttribute("aria-invalid");

        var errorMessage = document.getElementById("manaError");
        if (errorMessage) {
            errorMessage.parentNode.removeChild(errorMessage);
        }
    });

    document.body.addEventListener("BaseUpdated", function() {
        var errorMessage = document.getElementById("baseError");
        if (errorMessage) {
            errorMessage.parentNode.removeChild(errorMessage);
        }
    });

});

// theme switching
let atr = 'data-theme';
let localTheme = localStorage.getItem(atr)
if (localTheme) {
    setTheme(localTheme);
}

// if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches && localTheme == null) {
//     // dark mode
//     setTheme('dark');
// }

function toggleColorMode() {
    let el = document.documentElement;
    let atr_data = el.getAttribute(atr);

    if (atr_data == 'light') {
        setTheme('dark');
    } else {
        setTheme('light');
    }
}

function setTheme(theme) {
    let el = document.documentElement;
    el.getAttribute(atr);
    el.setAttribute(atr, theme);
    localStorage.setItem(atr, theme)
}

