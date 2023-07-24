import { DialogConfig } from 'kubeflow';

// --- Configs for the Confirm Dialogs ---
export function generateDeleteConfig(name: string): DialogConfig {
  return {
    title: $localize`Delete experiment ${name}?`,
    message: $localize`You cannot undo this action. Are you sure you want to delete this experiment?`,
    accept: $localize`DELETE`,
    confirmColor: 'warn',
    cancel: $localize`CANCEL`,
    error: '',
    applying: $localize`DELETING`,
    width: '600px',
  };
}
