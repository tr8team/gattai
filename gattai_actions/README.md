
<table>
<tr>
<td> File </td> <td> Fields </td><td>Description</td>
</tr>
<tr>
<td rowspan="3">
<b>gattai_actions/aws/get_elasticache_cluster_url.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/get_elasticache_cluster_url
vars:
  identifier: '"string"'
  profile: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The region where the resource is at</td>
</tr>
<tr>
<td>identifier<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="3">
<b>gattai_actions/aws/get_elasticache_node_grp_url.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/get_elasticache_node_grp_url
vars:
  identifier: '"string"'
  profile: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The region where the resource is at</td>
</tr>
<tr>
<td>identifier<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="3">
<b>gattai_actions/aws/get_elasticache_replica_grp_url.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/get_elasticache_replica_grp_url
vars:
  identifier: '"string"'
  profile: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The region where the resource is at</td>
</tr>
<tr>
<td>identifier<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="3">
<b>gattai_actions/aws/get_rds_cluster_url.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/get_rds_cluster_url
vars:
  identifier: '"string"'
  profile: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The region where the resource is at</td>
</tr>
<tr>
<td>identifier<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="3">
<b>gattai_actions/aws/get_rds_instance_url.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/get_rds_instance_url
vars:
  identifier: '"string"'
  profile: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The region where the resource is at</td>
</tr>
<tr>
<td>identifier<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="3">
<b>gattai_actions/aws/get_secgroup_id.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/get_secgroup_id
vars:
  group_name: '"string"'
  profile: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The region where the resource is at</td>
</tr>
<tr>
<td>group_name<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="4">
<b>gattai_actions/aws/get_secret_value.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/get_secret_value
vars:
  identifier: '"string"'
  profile: '"string"'
  property: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The region where the resource is at</td>
</tr>
<tr>
<td>identifier<br/><b>(required)</b></td>
<td>The secret identifier</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td>property<br/><i>(optional)</i></td>
<td>The property in secret to retrieve</td>
</tr>
<tr>
<td rowspan="4">
<b>gattai_actions/aws/query/elasticache_replica_grp.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/query/elasticache_replica_grp
vars:
  identifier: '"string"'
  profile: '"string"'
  query: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The region where the resource is at</td>
</tr>
<tr>
<td>identifier<br/><b>(required)</b></td>
<td>The elasticache replica grp identifier</td>
</tr>
<tr>
<td>query<br/><b>(required)</b></td>
<td>The elasticache replica grp query</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="4">
<b>gattai_actions/aws/query/rds_cluster.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/query/rds_cluster
vars:
  identifier: '"string"'
  profile: '"string"'
  query: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The regiion where the resource is at</td>
</tr>
<tr>
<td>identifier<br/><b>(required)</b></td>
<td>The rds cluster identifier</td>
</tr>
<tr>
<td>query<br/><b>(required)</b></td>
<td>The rds cluster query</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="4">
<b>gattai_actions/aws/query/rds_instance.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/query/rds_instance
vars:
  identifier: '"string"'
  profile: '"string"'
  query: '"string"'
  region: '"string"'

```

</td>
<td>region<br/><b>(required)</b></td>
<td>The regiion where the resource is at</td>
</tr>
<tr>
<td>identifier<br/><b>(required)</b></td>
<td>The rds instance identifier</td>
</tr>
<tr>
<td>query<br/><b>(required)</b></td>
<td>The rds instance query</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="4">
<b>gattai_actions/aws/query/security_grp.yaml</b>

```yaml
action: <repo_id>/gattai_actions/aws/query/security_grp
vars:
  filters:
    '"string"': '"string"'
  profile: '"string"'
  query: '"string"'
  region: '"string"'

```

</td>
<td>filters<br/><b>(required)</b></td>
<td>The security grp filters</td>
</tr>
<tr>
<td>query<br/><b>(required)</b></td>
<td>The security grp query</td>
</tr>
<tr>
<td>region<br/><b>(required)</b></td>
<td>The region where the resource is at</td>
</tr>
<tr>
<td>profile<br/><i>(optional)</i></td>
<td>The profile to use for accessing</td>
</tr>
<tr>
<td rowspan="5">
<b>gattai_actions/k8s/configmap.yaml</b>

```yaml
action: <repo_id>/gattai_actions/k8s/configmap
vars:
  fromEnvFile: '"string"'
  fromFile: '"string"'
  fromLiterals:
    '"string"': '"string"'
  name: '"string"'
  namespace: '"string"'

```

</td>
<td>name<br/><b>(required)</b></td>
<td>The kube configmap name</td>
</tr>
<tr>
<td>namespace<br/><b>(required)</b></td>
<td>The kube configmap namespace</td>
</tr>
<tr>
<td>fromFile<br/><i>(optional)</i></td>
<td>Create kube configmap from file</td>
</tr>
<tr>
<td>fromEnvFile<br/><i>(optional)</i></td>
<td>Create kube configmap from envfile</td>
</tr>
<tr>
<td>fromLiterals<br/><i>(optional)</i></td>
<td>Create kube configmap from literals</td>
</tr>
<tr>
<td rowspan="5">
<b>gattai_actions/k8s/secret_generic.yaml</b>

```yaml
action: <repo_id>/gattai_actions/k8s/secret_generic
vars:
  fromEnvFile: '"string"'
  fromFile: '"string"'
  fromLiterals:
    '"string"': '"string"'
  name: '"string"'
  namespace: '"string"'

```

</td>
<td>name<br/><b>(required)</b></td>
<td>The kube secret name</td>
</tr>
<tr>
<td>namespace<br/><b>(required)</b></td>
<td>The kube secret namespace</td>
</tr>
<tr>
<td>fromLiterals<br/><i>(optional)</i></td>
<td>Create kube secret from literals</td>
</tr>
<tr>
<td>fromFile<br/><i>(optional)</i></td>
<td>Create kube secret from file</td>
</tr>
<tr>
<td>fromEnvFile<br/><i>(optional)</i></td>
<td>Create kube secret from envfile</td>
</tr>
<tr>
<td rowspan="7">
<b>gattai_actions/k8s_rsc/external_secret.yaml</b>

```yaml
action: <repo_id>/gattai_actions/k8s_rsc/external_secret
vars:
  creationPolicy: '"string"'
  deletionPolicy: '"string"'
  k8sSecretName: '"string"'
  labels:
    '"string"': '"string"'
  name: '"string"'
  secretStoreKind: '"string"'
  secretStoreName: '"string"'

```

</td>
<td>name<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>secretStoreName<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>secretStoreKind<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>k8sSecretName<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>creationPolicy<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>deletionPolicy<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td>labels<br/><i>(optional)</i></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td rowspan="3">
<b>gattai_actions/op/item_get.yaml</b>

```yaml
action: <repo_id>/gattai_actions/op/item_get
vars:
  identifier: '"string"'
  label: '"string"'
  vault: '"string"'

```

</td>
<td>identifier<br/><b>(required)</b></td>
<td>The one password identifier</td>
</tr>
<tr>
<td>label<br/><b>(required)</b></td>
<td>The one password label</td>
</tr>
<tr>
<td>vault<br/><b>(required)</b></td>
<td>The one password vault</td>
</tr>
<tr>
<td rowspan="1">
<b>gattai_actions/print_output.yaml</b>

```yaml
action: <repo_id>/gattai_actions/print_output
vars:
  content: '"string"'

```

</td>
<td>content<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td rowspan="2">
<b>gattai_actions/terraform_rsc/write_tfvars.yaml</b>

```yaml
action: <repo_id>/gattai_actions/terraform_rsc/write_tfvars
vars:
  filename: '"string"'
  key_value_pairs:
    '"string"': '"string"'

```

</td>
<td>key_value_pairs<br/><b>(required)</b></td>
<td>The key value pairs for tfvars file</td>
</tr>
<tr>
<td>filename<br/><b>(required)</b></td>
<td>The filename to save the content into</td>
</tr>
<tr>
<td rowspan="2">
<b>gattai_actions/write_file.yaml</b>

```yaml
action: <repo_id>/gattai_actions/write_file
vars:
  content: '"string"'
  filename: '"string"'

```

</td>
<td>filename<br/><b>(required)</b></td>
<td>The filename to save the content into</td>
</tr>
<tr>
<td>content<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
<tr>
<td rowspan="2">
<b>gattai_actions/write_temp_file.yaml</b>

```yaml
action: <repo_id>/gattai_actions/write_temp_file
vars:
  content: '"string"'
  filename: '"string"'

```

</td>
<td>filename<br/><b>(required)</b></td>
<td>The filename to save the content into</td>
</tr>
<tr>
<td>content<br/><b>(required)</b></td>
<td>The content to be save if any</td>
</tr>
</table>
