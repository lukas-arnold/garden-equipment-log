Chart.defaults.elements.line.tension = 0.4;
Chart.defaults.interaction = {
    mode: "index",
    intersect: false
};

function formatValue(value, unit) {
    if (value == null) return "";

    const num = Number(value);
    if (Number.isNaN(num)) return value;

    switch (unit) {
        case "kg":
        case "cm":
            return (
                num.toLocaleString("de-DE", {
                    minimumFractionDigits: 1,
                    maximumFractionDigits: 1
                }) +
                " " +
                unit
            );

        case "%":
            return (
                num.toLocaleString("de-DE", {
                    maximumFractionDigits: 0
                }) +
                " " +
                unit
            );
        case "h/min":
            const totalMinutes = Math.floor(num);
            const hours = Math.floor(totalMinutes / 60);
            const minutes = totalMinutes % 60;

            if (hours > 0) {
                if (minutes === 0) {
                    return `${hours} h`;
                }

                return `${hours} h ${minutes} min`;
            }

            return `${minutes} min`;

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

    const formatXAxis = model.type === "bar"
        ? formatYear
        : formatDate;

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
                        callback: function(value) {
                            return formatXAxis(
                                this.getLabelForValue(value)
                            );
                        }
                    }
                },

                y: {
                    beginAtZero: true,

                ticks: {
                    stepSize: model.type === "bar" && unit === "h/min"
                        ? 60
                        : undefined,

                    callback(value) {
                        return formatValue(value, unit);
                    }
                }
                }
            },

            plugins: {
                tooltip: {
                    callbacks: {
                        title: (items) => {
                            return formatXAxis(items[0].label);
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
