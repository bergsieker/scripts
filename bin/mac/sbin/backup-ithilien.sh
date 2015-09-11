#!/bin/sh

# To restore:
#   rsync -v -a --exclude=/sbergsieker/Library/ sbergsieker@172.16.200.202::ithilien /Users 2>&1 > rsync.log
# If there are problems, check /mnt/HD_b2/logs/rsync.log. It is likely a router issue.

BACKUP_NAME=backup-$(date +%Y_%m_%d-%H_%M_%S)
LOG_NAME=${BACKUP_NAME}.txt
TMP_LOG_DIR=/tmp
FINAL_LOG_DIR=~/Documents/logs

# It is important that this NOT have a trailing slash.
SRC_DIR=/Users/sbergsieker
#SRC_DIR=/Users/sbergsieker/bin/zs/

RSYNC_TARGET_USER=sbergsieker
RSYNC_TARGET_ADDR=192.168.1.2
RSYNC_TARGET_NAME=ithilien
#RSYNC_TARGET_NAME=test

#DEBUG_OPTIONS='--dry-run -v'
#DEBUG_OPTIONS='-v'

echo Recall that this uses a different password than you would think...

EXCLUDES=
# The following files are excluded according to the advice at
# http://face.centosprime.com/macosxw/time-machine-default-exclusions/, which
# claims that TimeMachine excludes them by default. The command that site used
# to determine this is "sudo mdfind "com_apple_backup_excludeItem =
# 'com.apple.backupd'"".
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Application Support/MobileSync'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Application Support/SyncServices'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Caches'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Logs'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Mail/Envelope Index'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Mail/AvailableFeeds'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Mirrors'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/PubSub/Database'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/PubSub/Downloads'"
#EXCLUDES="${EXCLUDES} --exclude='sbergsieker/Library/PubSub/Feeds'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Safari/Icons.db'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Safari/HistoryIndex.sk'"

# This is excluded according to the advice at
# http://ryanblock.com/2008/05/good-folders-to-exclude-from-time-machine-backups/,
# which claims that opening attachments in Mail creates a copy of the
# attachment in this folder.
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Library/Mail Downloads'"

# Excluded because all music should live on the server already.
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Music/Amazon MP3'"
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Music/iTunes/iTunes Music'"

# Excluded because it shouldn't have anything useful.
#EXCLUDES="${EXCLUDES} --exclude='/sbergsieker/Downloads'"

# These are generic file names that we want to exclude.
#EXCLUDES="${EXCLUDES} --exclude='.DS_Store'"
#EXCLUDES="${EXCLUDES} --exclude='.Trash'"
#EXCLUDES="${EXCLUDES} --exclude='NoBackup'"


rsync \
  ${DEBUG_OPTIONS} \
  --log-file ${TMP_LOG_DIR}/${LOG_NAME} \
  -a \
  --delete \
  --delete-excluded \
  --exclude='/sbergsieker/Library/Application Support/MobileSync' \
  --exclude='/sbergsieker/Library/Application Support/SyncServices' \
  --exclude='/sbergsieker/Library/Caches' \
  --exclude='/sbergsieker/Library/Logs' \
  --exclude='/sbergsieker/Library/Mail/Envelope Index' \
  --exclude='/sbergsieker/Library/Mail/AvailableFeeds' \
  --exclude='/sbergsieker/Library/Mirrors' \
  --exclude='/sbergsieker/Library/PubSub/Database' \
  --exclude='/sbergsieker/Library/PubSub/Downloads' \
  --exclude='/sbergsieker/Library/PubSub/Feeds' \
  --exclude='/sbergsieker/Library/Safari/Icons.db' \
  --exclude='/sbergsieker/Library/Safari/HistoryIndex.sk' \
  --exclude='/sbergsieker/Library/Mail Downloads' \
  --exclude='/sbergsieker/Downloads' \
  --exclude='/sbergsieker/Music/Amazon MP3' \
  --exclude='/sbergsieker/Music/iTunes/iTunes Music' \
  --exclude='.DS_Store' \
  --exclude='.Trash' \
  --exclude='NoBackup' \
  ${SRC_DIR} \
  ${RSYNC_TARGET_USER}@${RSYNC_TARGET_ADDR}::${RSYNC_TARGET_NAME}

# When we run the backup command, we use a log file in a location that will not
# be backed up. Here, we copy the log file to the final location.
cp ${TMP_LOG_DIR}/${LOG_NAME} ${FINAL_LOG_DIR}/${LOG_NAME}

