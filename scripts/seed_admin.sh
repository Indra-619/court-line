#!/usr/bin/env bash
# Seed an admin user by email.
#
# Usage: ./scripts/seed_admin.sh user@example.com
#
# The email is passed to mongosh via an environment variable and read
# with process.env inside the eval script, so user input is never
# interpolated into JavaScript (no shell/JS injection surface).
set -euo pipefail

if [ $# -lt 1 ]; then
    echo "Usage: $0 <user-email>"
    exit 1
fi

EMAIL="$1"

# Sanity check: reject anything that is not a plausible email so the
# operator gets a clear error instead of a silent no-match.
if ! [[ "$EMAIL" =~ ^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$ ]]; then
    echo "Error: '$EMAIL' does not look like a valid email address." >&2
    exit 1
fi

MONGO_URI="${MONGO_URI:-mongodb://localhost:27017/booklapangan}"

SEED_EMAIL="$EMAIL" mongosh "$MONGO_URI" --quiet --eval '
const email = process.env.SEED_EMAIL;
const result = db.users.updateOne(
    { email: email },
    { $set: { role: "admin" } }
);
if (result.matchedCount === 0) {
    print("No user found with email " + email + ". The user must sign in with Google at least once before being promoted.");
    quit(1);
} else {
    print("User " + email + " is now an admin.");
}
'