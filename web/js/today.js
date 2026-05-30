document.addEventListener('DOMContentLoaded', function() {
    const today = new Date();
    const offset = today.getTimezoneOffset();
    const todayUTC = new Date(today.getTime() + offset * 60 * 1000);
    const currentDate = todayUTC.toISOString().split('T')[0];
    const localDateTime = new Date(today.getTime() - offset * 60 * 1000).toISOString().slice(0, 16);

    document.querySelectorAll('input[type="date"]').forEach(input => {
        if (!input.value) {
            input.value = currentDate;
        }
    });

    document.querySelectorAll('input[type="datetime-local"]').forEach(input => {
        if (!input.value) {
            input.value = localDateTime;
        }
    });
});
