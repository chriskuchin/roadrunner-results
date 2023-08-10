import { precacheAndRoute } from 'workbox-precaching/precacheAndRoute';
import { BackgroundSyncPlugin } from 'workbox-background-sync';
import { registerRoute } from 'workbox-routing';
import { NetworkOnly, NetworkFirst } from 'workbox-strategies';

precacheAndRoute(self.__WB_MANIFEST);

const bgSyncPlugin = new BackgroundSyncPlugin('failedWrites', {
  maxRetentionTime: 24 * 60, // Retry for max of 24 Hours (specified in minutes)
});

const statusPlugin = {
  fetchDidSucceed: ({ response }) => {
    if (response.status >= 500) {
      // Throwing anything here will trigger fetchDidFail.
      throw new Error('Server error.');
    }
    // If it's not 5xx, use the response as-is.
    return response;
  },
};

registerRoute(new RegExp(/\/api\/.*/), new NetworkFirst({}), 'GET')

registerRoute(
  new RegExp(/\/api\/.*/),
  new NetworkOnly({
    plugins: [statusPlugin, bgSyncPlugin],
  }),
  'POST'
);


self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});