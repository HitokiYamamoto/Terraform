/* NOTE:
    Secret Managerの無料枠は一つのリージョンに対して6個のバージョンまで
    無料枠の例:
        シークレットA : アクティブなバージョンが3つ
        シークレットB : アクティブなバージョンが2つ
        シークレットC : アクティブなバージョンが1つ
    有料になる例；
        シークレットA : アクティブなバージョンが4つ
        シークレットB : アクティブなバージョンが3つ

    減らす場合、"無効"ではなく"破棄"する必要がある。
    シークレットバージョンについては、
*/

resource "google_secret_manager_secret" "main" {
  secret_id = var.secret_id

  replication {
    auto {}
  }

  deletion_protection = false
}
