import { DialogConfig } from 'kubeflow';

// --- Configs for the Confirm Dialogs ---
export function getDeleteDialogConfig(
  name: string,
  namespace: string,
): DialogConfig {
  return {
    title: `Delete experiment`,
    message: `Are you sure you want to delete ${name} experiment from namespace ${namespace}?`,
    accept: 'DELETE',
    confirmColor: 'warn',
    cancel: 'CANCEL',
    error: '',
    applying: 'DELETING',
    width: '600px',
  };
}
