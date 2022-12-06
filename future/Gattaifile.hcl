// version = "v1"

// temp_folder = "./localtmp"

// enforce_targets {
//     develop = [
//         "mysql_password",
//         "return"
//     ]
//     #staging = enforce_targets.dev
//     #production = enforce_targets.dev
// }

// repo "local-src" {
//     src = "local"
//     config {
//         dir =  "gattai_actions"
//     }
// }

// repo "git-src" {
//     src = "git"
//     config {
//         url = "https://github.com/tr8team/gattai-libs.git"
//         branch = "develop"
//     }
// }

// target "develop" {
//     mysql_password {
//         action = "${repo.local-src.path}/op/item_get"
//         args {
//             identifier = "Okta YW local Testing"
//             label = "password"
//             vault = "Private"
//         }
//     }
//     create_envfile {
//         action = "${repo.local-src.path}/write_temp_file"
//         args {
//         filename = ".env.remote-testing"
//         content = <<EOF
//             DB_HOST=engineering-remote-testing-transactional-database
//             DB_USERNAME=admin
//             DB_PASSWORD=${fetch(target.develop.mysql_password)}
//             REDIS_HOST=engineering-remote-testing-cache-master
//             EOF
//         }
//     }
//     return {
//       action = "${repo.local-src.path}/k8s/configmap"
//       args {
//         name = "envfile"
//         namespace = "default"
//         fromEnvFile = fetch(target.develop.create_envfile)
//       }
//     }
// }
