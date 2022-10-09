<table>
<tr>
<td> File </td> <td> Fields </td><td>Description</td>
</tr>
<tr>
<td rowspan="5">
<b>k8s/configmap.yaml</b>

```yaml
action: <repo_id>/k8s/configmap
vars:
  fromEnvFile: '"string"'
  fromFile: '"string"'
  fromLiterals:
    '"string"': '"string"'
  name: '"string"'
  namespace: '"string"'
```

</td>
<td>fromEnvFile<br/><i>(optional)</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>fromFile<br/><i>(optional)</i></td>
<td>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</td>
</tr>
<tr>
<td>fromLiterals<br/><i>(optional)</i></td>
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
