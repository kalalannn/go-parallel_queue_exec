package messages

const StartShutdownServer = "Gracefully shutdown Server..."
const StartShutdownFiber = "Shutting down Fiber..."
const EndShutdownFiber = "Fiber shutdown finished."

const StartShutdownExecutorService = "Try to shutdown executor service..."
const EndShutdownExecutorServiceSuccess = "Executor service shutdown completly finished. Cool!"
const EndShutdownExecutorServiceTimeout = "Executor service shutdown INCOMPLETLY finished. Don't care."
