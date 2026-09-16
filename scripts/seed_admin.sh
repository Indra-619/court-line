#!/usr/bin/env bash
# Seed an admin user by email.
#
# Usage: ./scripts/seed_admin.sh user@example.com
set -euo pipefail

if [ $# -lt 1 ]; then
    echo "Usage: $0 <user-email>"
    exit 1
fi

EMAIL="$1"

mongosh "mongodb://localhost:27017/booklapangan" --quiet --eval "
const result = db.users.updateOne(
    { email: '${EMAIL}' },
    { \$set: { role: 'admin' } }
);
if (result.matchedCount === 0) {
    print('No user found with email ${EMAIL}. The user must sign in with Google at least once before being promoted.');
    quit(1);
} else {
    print('User ${EMAIL} is now an admin.');
}
"
