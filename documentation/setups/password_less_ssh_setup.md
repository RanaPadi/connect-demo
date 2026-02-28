# SSH Authentication with Private/Public Keys

To configure SSH authentication to work with private/public keys on Ubuntu, follow these steps:

1. Generate SSH key pair:
   - Open a terminal.
   - Run the following command to generate a new SSH key pair:

     ```bash
     ssh-keygen -t ed25519 -C SSH_DN-CN -f ~/.ssh/id_ed25519_SSH_DN-CN -N ''
     ```

2. Copy the public key to the remote server:
   - Run the following command to copy the public key to the remote server:

     ```bash
     ssh-copy-id ~/.ssh/id_ed25519_SSH_DN-CN username@remote_host
     ```

   - Enter the password for the remote server when prompted.

3. Test the SSH connection:
   - Run the following command to test the SSH connection using the key pair:

     ```bash
     ssh username@remote_host
     ```

   - If the connection is successful, you should be logged in to the remote server without entering a password.

4. Disable password authentication:
   - Open the SSH server configuration file using a text editor:

     ```bash
     sudo nano /etc/ssh/sshd_config
     ```

   - Add the line that says `PubkeyAuthentication yes`
   - Add the line that says `ChallengeResponseAuthentication no`
   - Find the line that says `#PasswordAuthentication yes` and change it to `PasswordAuthentication no`.
   - Find the line that says `#UsePAM yes` and change it to `UsePAM no`.
   - Save the file and exit the text editor.
   - Restart the SSH service for the changes to take effect:

     ```bash
     sudo systemctl restart sshd
     ```

That's it! You have successfully configured SSH authentication to work with private/public keys on Ubuntu.
