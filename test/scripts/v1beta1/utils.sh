_activate_service_account() {

  # Sometimes activating SA can fail, we try to restart it 5 times.
  _count_attempts=1
  _max_attempts=5
  _is_activated=false

  while [[ $_is_activated=false && $_count_attempts -le $_max_attempts ]]; do
    echo "Activating service-account"
    gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS} && $_is_activated=true || count_attempts=$((count_attempts + 1))

    if [ $_is_activated=false ]; then
      echo "gcloud activate service account failed, restart"
      sleep 1
    fi

    _count_attempts=$((_count_attempts + 1))
  done

  # If account was not activated exit the script
  if [ $_is_activated=false ]; then
    echo "Unable to activate gcloud service account!"
    exit 0
  fi

}
