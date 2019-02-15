# Tips

1. MPI can spawn different number of copies. It is controlled by mpirun -n inside the Dockerfile.

    ```
        mpirun -n <number_of_copies>
    ```
    
    **Note.** Each copy will utilise 1 CPU. You can binding each process a CPU using `-cpu-slot`. For more reference visit [mpirun docuentation](https://www.open-mpi.org/doc/v3.0/man1/mpirun.1.php).
   