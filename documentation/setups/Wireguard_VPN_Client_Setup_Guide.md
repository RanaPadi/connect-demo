# Client Side Setup

## Windows

1. **Download and Install WireGuard:**
   - Visit the [official WireGuard website](https://www.wireguard.com/install/) and download the Windows installer.
   - Run the installer and follow the prompts to complete the installation.

2. **Generate a Configuration File:**
   - Open the WireGuard application.
   - Click on “Add Tunnel” and then select “Add empty tunnel…”.
   - The application will generate a public and private key automatically. Copy these keys to use in the server configuration if needed.

3. **Create a Configuration File:**
   - You can either create the configuration file manually or get it from your VPN provider.
   - A typical configuration file (`wg0.conf`) looks like this:

     ```shell
     [Interface]
     PrivateKey = YOUR_PRIVATE_KEY
     Address = 10.0.0.X/32
     DNS = OpenWRT_ROUTER_GATEWAY_IP

     [Peer]
     PublicKey = SERVER_PUBLIC_KEY
     Endpoint = YOUR_SERVER_IP:51820
     AllowedIPs = 0.0.0.0/0
     ```

   - Replace `YOUR_PRIVATE_KEY`, `SERVER_PUBLIC_KEY`, `10.0.0.X`, `OpenWRT_ROUTER_GATEWAY_IP`, and `YOUR_SERVER_IP` with the actual values.

4. **Import the Configuration:**
   - In the WireGuard application, click on “Import Tunnel(s) from file…”.
   - Select your configuration file: `wg0.conf`, and it will be added to the application.

5. **Activate the VPN:**
   - Toggle the switch to activate the VPN. Your connection should now be active.

## macOS

1. **Download and Install WireGuard:**
   - Go to the [App Store](https://apps.apple.com/us/app/wireguard/id1451685025?mt=12) and search for "WireGuard".
   - Install the WireGuard application.

2. **Generate a Configuration File:**
   - Open the WireGuard application.
   - Click on “Create Tunnel” and then “Generate Keypair”.
   - This will generate a public and private key pair. Save these keys for server configuration.

3. **Create a Configuration File:**
   - Similar to Windows, the configuration file will look like:

     ```shell
     [Interface]
     PrivateKey = YOUR_PRIVATE_KEY
     Address = 10.0.0.X/32
     DNS = OpenWRT_ROUTER_GATEWAY_IP

     [Peer]
     PublicKey = SERVER_PUBLIC_KEY
     Endpoint = YOUR_SERVER_IP:51820
     AllowedIPs = 0.0.0.0/0
     ```

   - Replace `YOUR_PRIVATE_KEY`, `SERVER_PUBLIC_KEY`, `10.0.0.X`, `OpenWRT_ROUTER_GATEWAY_IP`, and `YOUR_SERVER_IP` with the actual values.

4. **Import the Configuration:**
   - In the WireGuard app, click on “Import Tunnel(s) from file…”.
   - Select your configuration file: `wg0.conf`.

5. **Activate the VPN:**
   - Toggle the switch next to the tunnel name to connect.

## Ubuntu/Debian

1. **Install WireGuard:**
   - Update your package list and install WireGuard:

     ```shell
     sudo apt update
     sudo apt install wireguard
     ```

2. **Generate Keys:**
   - Generate your private and public keys:

     ```shell
     umask 077
     wg genkey | tee privatekey | wg pubkey > publickey
     ```

   - Save these keys in a secure location.

3. **Create a Configuration File:**
   - Create the configuration file at `/etc/wireguard/wg0.conf`:

     ```shell
     sudo nano /etc/wireguard/wg0.conf
     ```

   - Populate it with:

     ```shell
     [Interface]
     PrivateKey = YOUR_PRIVATE_KEY
     Address = 10.0.0.X/32
     DNS = OpenWRT_ROUTER_GATEWAY_IP

     [Peer]
     PublicKey = SERVER_PUBLIC_KEY
     Endpoint = YOUR_SERVER_IP:51820
     AllowedIPs = 0.0.0.0/0
     ```

- Replace `YOUR_PRIVATE_KEY`, `SERVER_PUBLIC_KEY`, `10.0.0.X`, `OpenWRT_ROUTER_GATEWAY_IP`, and `YOUR_SERVER_IP` with the actual values.

## Start and Enable WireGuard:**

- If you did not create the Wireguard configuration file: `wg0.conf` in the previous step, instead you have the said configuration file from your admin, then copy it to the directory: '/etc/wireguargd.`
  - Start the WireGuard interface:

     ```shell
     sudo wg-quick up wg0
     ```

  - Enable the service to start on boot:

     ```shell
     sudo systemctl enable wg-quick@wg0
     ```

1. **Verify the Connection:**
   - You can check the status with:

     ```shell
     sudo wg
     ```

2. **Troubleshooting:**
   - If the connection doesn't work, check logs:

     ```shell
     sudo journalctl -xe
     ```
