Chart.defaults.elements.line.tension = 0.4;
Chart.defaults.interaction = {
    mode: "index",
    intersect: false
};

function isIsoDate(value) {
    return typeof value === "string" &&
        /^\d{4}-\d{2}-\d{2}$/.test(value);
}

function formatValue(value, unit) {
    if (value == null) return "";

    const num = Number(value);
    if (Number.isNaN(num)) return value;

    switch (unit) {
        case "kg":
        case "cm":
        case "h":
            return (
                num.toLocaleString("de-DE", {
                    minimumFractionDigits: 1,
                    maximumFractionDigits: 1
                }) +
                " " +
                unit
            );

        case "min":
        case "%":
            return (
                num.toLocaleString("de-DE", {
                    maximumFractionDigits: 0
                }) +
                " " +
                unit
            );

        default:
            return num.toLocaleString("de-DE");
    }
}

function renderChart(canvasId, model) {
    const ctx = document.getElementById(canvasId);

    const datasets = model.sets.map(set => ({
        label: set.label,
        data: set.data
    }));

    const unit = model.sets.length > 0
        ? model.sets[0].unit
        : "";

    new Chart(ctx, {
        type: model.type || "line",

        data: {
            labels: model.labels,
            datasets
        },

        options: {
            responsive: true,

            scales: {
                x: {
                    ticks: {
                        callback(value) {

                            const label =
                                this.getLabelForValue(value);

                            return isIsoDate(label)
                                ? formatDate(label)
                                : label;
                        }
                    }
                },

                y: {
                    beginAtZero: true,

                    ticks: {
                        callback(value) {
                            return formatValue(
                                value,
                                unit
                            );
                        }
                    }
                }
            },

            plugins: {
                tooltip: {
                    callbacks: {

                        title(items) {

                            const label = items[0].label;

                            return isIsoDate(label)
                                ? formatDate(label)
                                : label;
                        },

                        label(ctx) {

                            const unit =
                                model.sets[ctx.datasetIndex]?.unit || "";

                            return (
                                ctx.dataset.label +
                                ": " +
                                formatValue(ctx.raw, unit)
                            );
                        }
                    }
                }
            }
        }
    });
}