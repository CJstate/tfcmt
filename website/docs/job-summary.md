---
sidebar_position: 575
---

# Output the result to the GitHub Actions job summary

[#1367](https://github.com/suzuki-shunsuke/tfcmt/issues/1367)

tfcmt can append the rendered result of `terraform plan` and `terraform apply` to the
[GitHub Actions job summary](https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#adding-a-job-summary)
in addition to posting a pull request comment.

This feature is opt-in because the job summary is only available in GitHub Actions.

## Enable

Use the `--job-summary` option:

```sh
tfcmt --job-summary plan -- terraform plan
tfcmt --job-summary apply -- terraform apply
```

Or set the environment variable `TFCMT_JOB_SUMMARY`:

```sh
TFCMT_JOB_SUMMARY=true tfcmt plan -- terraform plan
```

Or enable it in the configuration file:

```yaml
job_summary: true
```

## Requirements

The environment variable `GITHUB_STEP_SUMMARY` must be set to the path of the job summary file.
GitHub Actions sets this environment variable automatically.
tfcmt appends the result to the file instead of overwriting it, so results of multiple steps are preserved.

## 1 MiB limit

GitHub Actions limits a job summary to 1 MiB per step.
If the result would exceed the limit, tfcmt truncates it and appends a notice.
Content already written by other steps is preserved.

## Notes

`--job-summary` can be used with `--output`. In that case tfcmt writes the result to both the local file and the job summary.
