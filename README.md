# mysql-vcs
VCS for MySQL Database - git style commit, log, checkout of the schema+data

Highlights
----------

* `commit` snapshot of database - structure, data, triggers and stored routines.
* Display `log` of comitted versions
* `checkout` any commit (restoring backups)
* Update latest commit using `commit --amend`
* To keep snopshot data concies and minimul
    * Easy `delete` of too older or unnecessary commits
    * Snapshots stored in compressed bytes (using DEFLATE)

    Technically, it's not a VCS. 
    It's just a simple tool to help developers to keep histry of 
    database backups and restore them with ease.

