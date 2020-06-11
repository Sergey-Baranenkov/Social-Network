export default function Debounce(f, t) {
    let lastCall = 0;
    let lastCallTimer;
    return function (...args) {
        let previousCall = lastCall;
        lastCall = Date.now();
        if (previousCall && ((lastCall- previousCall) <= t)) {
            clearTimeout(lastCallTimer);
        }
        lastCallTimer = setTimeout(() => f(...args), t);
    }
}