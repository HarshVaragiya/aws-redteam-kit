# AWS RedTeam Kit

Aka - Aws Ransomware Simulation Kit

This repository aims to be a PoC for the research project "Perfecting Ransomware on AWS."

# Note
- DO NOT ENCRYPT ANY DATA WITH THIS THAT YOU CAN'T AFFORD TO LOSE
- This project is licensed under [MIT](LICENSE)
- The Author is not responsible for any data loss occouring due to this project
- Docker containers are ephemeral
- This project is just a proof-of-concept ransomware simulation toolkit. It is NOT PRODUCTION GRADE.

# KeySwitch

The tool simulates a ransomware attack on a target by encrypting EBS volumes with the specified key.

**Use this tool only if you know exactly what you are doing. It deletes all snapshots in an AWS Region**

The tool works in the following manner :
1. Generate list of all EBS volumes in the region
2. Take snapshots of all volumes
3. Wait for snapshots to become Available.
4. Create new Volumes using the snapshots and specify encryption configuration as the KMS key.
5. Delete all the snapshots in the AWS Region (Nuke).

# Simulating a Ransomware Attacking using ARK
https://medium.com/@harsh8v/redefining-ransomware-attacks-on-aws-using-aws-kms-xks-dea668633802
