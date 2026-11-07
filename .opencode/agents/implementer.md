---
description: Implements PROMPT.md during an automated CI run
mode: primary
permissions:
  - action: edit
    resource: ".github/**"
    effect: deny
  - action: edit
    resource: ".opencode/**"
    effect: deny
  - action: edit
    resource: "PROMPT.md"
    effect: deny
  - action: edit
    resource: "opencode.json"
    effect: deny
  - action: edit
    resource: "README.md"
    effect: deny
  - action: shell
    resource: "git commit *"
    effect: deny
  - action: shell
    resource: "git push *"
    effect: deny
  - action: shell
    resource: "gh *"
    effect: deny
  - action: webfetch
    resource: "*"
    effect: deny
  - action: websearch
    resource: "*"
    effect: deny
---
