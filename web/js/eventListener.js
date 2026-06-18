document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll(".date")
        .forEach(el => el.textContent = formatDate(el.textContent));

    document.querySelectorAll(".datetime")
        .forEach(el => el.textContent = formatDateTime(el.textContent));

    document.querySelectorAll(".bottle-purchase-date")
        .forEach(el => el.textContent = formatBottlePurchaseDate(el.textContent));

    const now = new Date();
    const offset = now.getTimezoneOffset() * 60000;

    const currentDate = new Date(now.getTime() + offset)
        .toISOString()
        .split('T')[0];

    const localDateTime = new Date(now.getTime() - offset)
        .toISOString()
        .slice(0, 16);

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