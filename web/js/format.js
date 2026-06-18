const formatDate = date =>
    new Date(date).toLocaleDateString("default", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit"
    });

const formatDateTime = date =>
    new Date(date).toLocaleDateString("default", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit"
    });

const formatBottlePurchaseDate = date =>
    new Date(date).toLocaleDateString("default", {
        year: "numeric",
        month: "2-digit",
    });
