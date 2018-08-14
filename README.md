GPU Exporter
=========================

This is a Prometheus Exporter for Collect Nvidia GPU metrics. It uses the Nvidia [NVML Go Bindings](https://github.com/NVIDIA/gpu-monitoring-tools), not call the `nvidia-smi` binary.


## Running

`./gpu_exporter run -v 5 --logtostderr`

```
# HELP nvidia_gpu_memory_total_mb Memory Total of the GPU device in MB
# TYPE nvidia_gpu_memory_total_mb gauge
nvidia_gpu_memory_total_mb{name="Tesla P4",path="/dev/nvidia0",uuid="GPU-9bb57e8e-f94d-4b95-7c33-ee279c8fb75c"} 7606
nvidia_gpu_memory_total_mb{name="Tesla P4",path="/dev/nvidia1",uuid="GPU-2d9067e4-8fab-5820-6228-69e2e74b3d58"} 7606

# HELP nvidia_gpu_memory_used_mb Memory used by the GPU device in MB
# TYPE nvidia_gpu_memory_used_mb gauge
nvidia_gpu_memory_used_mb{name="Tesla P4",path="/dev/nvidia0",uuid="GPU-9bb57e8e-f94d-4b95-7c33-ee279c8fb75c"} 0
nvidia_gpu_memory_used_mb{name="Tesla P4",path="/dev/nvidia1",uuid="GPU-2d9067e4-8fab-5820-6228-69e2e74b3d58"} 0

# HELP nvidia_gpu_num_devices Number of Nvidia GPU devices
# TYPE nvidia_gpu_num_devices gauge
nvidia_gpu_num_devices 2
# HELP nvidia_gpu_power_usage Power usage of the GPU device in watts
# TYPE nvidia_gpu_power_usage gauge
nvidia_gpu_power_usage{name="Tesla P4",path="/dev/nvidia0",uuid="GPU-9bb57e8e-f94d-4b95-7c33-ee279c8fb75c"} 7
nvidia_gpu_power_usage{name="Tesla P4",path="/dev/nvidia1",uuid="GPU-2d9067e4-8fab-5820-6228-69e2e74b3d58"} 6

# HELP nvidia_gpu_temperature Temperature of the GPU device in celsius
# TYPE nvidia_gpu_temperature gauge
nvidia_gpu_temperature{name="Tesla P4",path="/dev/nvidia0",uuid="GPU-9bb57e8e-f94d-4b95-7c33-ee279c8fb75c"} 45
nvidia_gpu_temperature{name="Tesla P4",path="/dev/nvidia1",uuid="GPU-2d9067e4-8fab-5820-6228-69e2e74b3d58"} 39
```