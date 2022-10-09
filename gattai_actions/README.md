<table>
<tr>
<td> File </td> <td> How to Use </td><td> Fields </td><td></td><td></td>
</tr>
<tr>
<td rowspan="5"> k8s/configmap.yaml </td>
<td rowspan="5">

```yaml
action: <repo_id>/k8s/configmap
vars:
  from_env_file: '"string"'
  from_file: '"string"'
  from_literals:
    '"string"': '"string"'
  name: '"string"'
  namespace: '"string"'
```

</td>
<td>from_env_file</td>
<td><i>optional</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>from_file</td>
<td><i>optional</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>from_literals</td>
<td><i>optional</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>name</td>
<td><b>required</b></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>namespace</td>
<td><b>required</b></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td rowspan="3"> k8s_rsc/external_secret.yaml </td>
<td rowspan="3">

```yaml
action: <repo_id>/k8s_rsc/external_secret
vars:
  data:
    - {}
  name: '"string"'
  secretStoreName: '"string"'
```

</td>
<td>name</td>
<td><b>required</b></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>secretStoreName</td>
<td><b>required</b></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>data</td>
<td><i>optional</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
