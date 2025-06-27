#!/usr/bin/env python3


import yaml
import sys
import json
from typing import Dict, List, Tuple, Any
import re

def load_manifests(file_path: str) -> List[Dict]:
    """Load YAML manifests from file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    docs = []
    for doc_str in content.split('---'):
        doc_str = doc_str.strip()
        if doc_str:
            try:
                doc = yaml.safe_load(doc_str)
                if doc:  
                    docs.append(doc)
            except yaml.YAMLError as e:
                print(f"Error parsing YAML: {e}")
                continue
    
    return docs

def clean_helm_metadata(obj: Any) -> Any:
    """Remove Helm-specific metadata that should be ignored in comparison."""
    if isinstance(obj, dict):
        cleaned = {}
        for key, value in obj.items():
            if key == "metadata" and isinstance(value, dict):
                # Clean metadata section
                cleaned_metadata = {}
                for meta_key, meta_value in value.items():
                    if meta_key == "labels" and isinstance(meta_value, dict):
                        # Remove only Helm-specific labels
                        cleaned_labels = {}
                        for label_key, label_value in meta_value.items():
                            if not label_key.startswith(('helm.sh/', 'app.kubernetes.io/')):
                                cleaned_labels[label_key] = label_value
                        if cleaned_labels:  # Only add if there are remaining labels
                            cleaned_metadata[meta_key] = cleaned_labels
                    elif meta_key == "annotations" and isinstance(meta_value, dict):
                        # Remove only Helm-specific annotations
                        cleaned_annotations = {}
                        for ann_key, ann_value in meta_value.items():
                            if not ann_key.startswith(('helm.sh/', 'meta.helm.sh/')):
                                cleaned_annotations[ann_key] = ann_value
                        if cleaned_annotations:  # Only add if there are remaining annotations
                            cleaned_metadata[meta_key] = cleaned_annotations
                    else:
                        # Keep all other metadata fields as-is
                        cleaned_metadata[meta_key] = meta_value
                cleaned[key] = cleaned_metadata
            else:
                cleaned[key] = clean_helm_metadata(value)
        return cleaned
    elif isinstance(obj, list):
        return [clean_helm_metadata(item) for item in obj]
    else:
        return obj

def normalize_kustomize_refs(obj: Any, path: str = "") -> Any:
    """Normalize Kustomize hash suffixes in secret/configmap references throughout the manifest."""
    if isinstance(obj, dict):
        normalized = {}
        for key, value in obj.items():
            current_path = f"{path}.{key}" if path else key
            
            # Normalize secret/configmap references in common locations
            if key == "name" and isinstance(value, str):
                if any(ref_pattern in path for ref_pattern in [
                    'secretKeyRef', 'configMapKeyRef', 'configMapRef', 'secret', 'volumes'
                ]):
                    value = re.sub(r'-[a-z0-9]{10}$', '', value)
                elif 'volumes' in path and key == 'secretName':
                    value = re.sub(r'-[a-z0-9]{10}$', '', value)
                elif 'volumes' in path and 'configMap' in path:
                    value = re.sub(r'-[a-z0-9]{10}$', '', value)
            
            normalized[key] = normalize_kustomize_refs(value, current_path)
        return normalized
    elif isinstance(obj, list):
        return [normalize_kustomize_refs(item, path) for item in obj]
    else:
        return obj

def normalize_manifest(manifest: Dict) -> Dict:
    """Normalize manifest by removing/standardizing certain fields."""
    normalized = manifest.copy()
    
    # Clean Helm-specific metadata
    normalized = clean_helm_metadata(normalized)
    
    # Normalize Kustomize hash references
    normalized = normalize_kustomize_refs(normalized)
    
    if 'metadata' in normalized and 'name' in normalized['metadata']:
        kind = normalized.get('kind', '')
        if kind in ['Secret', 'ConfigMap']:
            name = normalized['metadata']['name']
            normalized['metadata']['name'] = re.sub(r'-[a-z0-9]{10}$', '', name)
    
    if 'metadata' in normalized:
        metadata = normalized['metadata']
        
        metadata.pop('generation', None)
        
        metadata.pop('resourceVersion', None)
        
        metadata.pop('uid', None)
        
        metadata.pop('creationTimestamp', None)
        
        metadata.pop('managedFields', None)
    
    normalized.pop('status', None)
    
    def remove_empty_values(obj):
        if isinstance(obj, dict):
            return {k: remove_empty_values(v) for k, v in obj.items() 
                   if v is not None and v != {} and v != []}
        elif isinstance(obj, list):
            return [remove_empty_values(item) for item in obj if item is not None]
        else:
            return obj
    
    return remove_empty_values(normalized)

def get_resource_key(manifest: Dict) -> str:
    """Generate a unique key for the resource."""
    kind = manifest.get('kind', 'Unknown')
    name = manifest.get('metadata', {}).get('name', 'unknown')
    namespace = manifest.get('metadata', {}).get('namespace', '')
    
    if kind in ['Secret', 'ConfigMap']:
        name = re.sub(r'-[a-z0-9]{10}$', '', name)
    
    if namespace:
        return f"{kind}/{namespace}/{name}"
    else:
        return f"{kind}/{name}"

def deep_diff(obj1: Any, obj2: Any, path: str = "") -> List[str]:
    """Compare two objects and return list of differences."""
    differences = []
    
    if type(obj1) != type(obj2):
        differences.append(f"{path}: type mismatch ({type(obj1).__name__} vs {type(obj2).__name__})")
        return differences
    
    if isinstance(obj1, dict):
        all_keys = set(obj1.keys()) | set(obj2.keys())
        for key in sorted(all_keys):
            key_path = f"{path}.{key}" if path else key
            if key not in obj1:
                differences.append(f"{key_path}: missing in kustomize")
            elif key not in obj2:
                differences.append(f"{key_path}: missing in helm")
            else:
                differences.extend(deep_diff(obj1[key], obj2[key], key_path))
    
    elif isinstance(obj1, list):
        if len(obj1) != len(obj2):
            differences.append(f"{path}: list length mismatch ({len(obj1)} vs {len(obj2)})")
        else:
            for i, (item1, item2) in enumerate(zip(obj1, obj2)):
                differences.extend(deep_diff(item1, item2, f"{path}[{i}]"))
    
    elif obj1 != obj2:
        differences.append(f"{path}: '{obj1}' != '{obj2}'")
    
    return differences

def compare_manifests(kustomize_file: str, helm_file: str, scenario: str) -> bool:
    """Compare Kustomize and Helm manifests."""
    kustomize_manifests = load_manifests(kustomize_file)
    helm_manifests = load_manifests(helm_file)
    
    kustomize_resources = {}
    helm_resources = {}
    
    for manifest in kustomize_manifests:
        normalized = normalize_manifest(manifest)
        key = get_resource_key(normalized)
        kustomize_resources[key] = normalized
    
    for manifest in helm_manifests:
        normalized = normalize_manifest(manifest)
        key = get_resource_key(normalized)
        helm_resources[key] = normalized
    
    kustomize_keys = set(kustomize_resources.keys())
    helm_keys = set(helm_resources.keys())
    
    common_keys = kustomize_keys & helm_keys
    only_in_kustomize = kustomize_keys - helm_keys
    only_in_helm = helm_keys - kustomize_keys
    
    expected_helm_extras = {
        'Secret/kubeflow/katib-webhook-cert',  # Webhook certificates
    }
    
    unexpected_helm_extras = only_in_helm - expected_helm_extras
    
    differences_found = []
    success = True
    
    if only_in_kustomize:
        print(f"Resources only in Kustomize:")
        for key in sorted(only_in_kustomize):
            print(f"  - {key}")
        success = False
        differences_found.extend(only_in_kustomize)
    
    if unexpected_helm_extras:
        print(f"Unexpected resources only in Helm:")
        for key in sorted(unexpected_helm_extras):
            print(f"  - {key}")
        success = False
        differences_found.extend(unexpected_helm_extras)
    
    # Compare common resources
    for key in sorted(common_keys):
        kustomize_resource = kustomize_resources[key]
        helm_resource = helm_resources[key]
        
        differences = deep_diff(kustomize_resource, helm_resource)
        
        if differences:
            print(f"Differences in {key}:")
            differences_found.append(key)
            success = False
            
            for diff in differences[:10]:
                print(f"   {diff}")
            if len(differences) > 10:
                print(f"   ... and {len(differences) - 10} more differences")
    
    if not success:
        print(f"Found differences in {len(differences_found)} resources")
        return False
    
    return True

if __name__ == "__main__":
    if len(sys.argv) < 4 or len(sys.argv) > 6:
        print("Usage: python compare_manifests.py <kustomize_file> <helm_file> <scenario> [namespace] [--verbose]")
        sys.exit(1)
    
    kustomize_file = sys.argv[1]
    helm_file = sys.argv[2]
    scenario = sys.argv[3]
    
    success = compare_manifests(kustomize_file, helm_file, scenario)
    sys.exit(0 if success else 1) 