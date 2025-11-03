export const toast = {
  success(message: string) {
    if (typeof window !== "undefined") {
      // eslint-disable-next-line no-console
      console.info(`[success] ${message}`);
    }
  },
  error(message: string) {
    if (typeof window !== "undefined") {
      // eslint-disable-next-line no-console
      console.error(`[error] ${message}`);
    }
  },
};
