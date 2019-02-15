# Tips

1. NVIDIA GPUs can now be via container level resource requirements using the resource name nvidia.com/gpu:
    ```
      resources:
        limits:
            nvidia.com/gpu: 2 # requesting 2 
    ```
    **Keep in mind!** The number of GPUs used by workers and master should be less of equal to the number of available GPUs on your cluster/system. If you should have less, then we recommend you to reduce the number of workers, or use master ony (in case you have 1 GPU).
    
2. MPI can spawn different number of copies. It is controlled by mpirun -n inside the Dockerfile.

    ```
        mpirun -n <number_of_copies>
    ```
    
    **Note.** Each copy will utilise 1 GPU.
   