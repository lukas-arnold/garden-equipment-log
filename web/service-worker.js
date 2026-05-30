const CACHE = "garden-equipment-log";

self.addEventListener("install", event => {
  self.skipWaiting();
});

self.addEventListener("activate", event => {
  event.waitUntil(
    (async () => {
      const keys = await caches.keys();

      await Promise.all(
        keys
          .filter(key => key !== CACHE)
          .map(key => caches.delete(key))
      );

      await clients.claim();
    })()
  );
});

self.addEventListener("fetch", event => {
  const req = event.request;
  const url = new URL(req.url);

  // only cache GET requests
  if (req.method !== "GET") return;

  // only cache http/https
  if (url.protocol !== "http:" && url.protocol !== "https:") return;

  event.respondWith(
    (async () => {
      try {
        const res = await fetch(req);

        // only cache successful responses
        if (res && res.status === 200) {
          const cache = await caches.open(CACHE);
          await cache.put(req, res.clone());
        }

        return res;
      } catch (err) {
        const cached = await caches.match(req);
        return cached || new Response("Offline", {
          status: 503,
          statusText: "Offline",
        });
      }
    })()
  );
});
