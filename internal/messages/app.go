package messages

const StartShutdownFiber = "Gracefully shutdown Fiber..."
const EndShutdownFiber = "Fiber shutdown finished."

const StartShutdownExecutorService = "Try to shutdown executor service..."
const EndShutdownExecutorServiceSuccess = "Executor service shutdown completly finished. Cool!"
const EndShutdownExecutorServiceTimeout = "Executor service shutdown INCOMPLETLY finished. Don't care."

const WaitForExecutorShutdownWithTimeout = "Waiting for executor shutdown (timeout: %s)."
const ExecutorShutdownFinished = "Executor shutdown finished."
const ExecutorShutdownTimeout = "Executor shutdown timeout."
