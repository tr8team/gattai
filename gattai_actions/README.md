<table>
<tr>
<td> File </td> <td> Fields </td><td></td>
</tr>
<tr>
<td rowspan="5">
<b>k8s/configmap.yaml</b>

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
<td>from_env_file<br/><i>(optional)</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>from_file<br/><i>(optional)</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>from_literals<br/><i>(optional)</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>name<br/><b>(required)</b></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>namespace<br/><b>(required)</b></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td rowspan="3">
<b>k8s_rsc/external_secret.yaml</b>

```yaml
action: <repo_id>/k8s_rsc/external_secret
vars:
  data:
    - {}
  name: '"string"'
  secretStoreName: '"string"'
```

</td>
<td>name<br/><b>(required)</b></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>secretStoreName<br/><b>(required)</b></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>data<br/><i>(optional)</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
