# WireGuard VPN Server Setup on OpenWRT Router

## Prerequisites

- An OpenWRT router with sufficient resources.
- Basic knowledge of SSH and OpenWRT's web interface (LuCI).
- Access to the router's web interface or SSH.

## Step 1: Install WireGuard on OpenWRT

1. **Access the Router:**
   - Log in to your OpenWRT router via the web interface or SSH.

2. **Update Package Lists:**
   - Update the package lists to ensure you have the latest versions available:

     ```shell
     opkg update
     ```

3. **Install WireGuard Packages:**
   - Install WireGuard packages using the following command:

     ```shell
     opkg install wireguard wireguard-tools luci-proto-wireguard
     ```

## Step 2: Generate WireGuard Keys

1. **Generate Server Keys:**
   - You need to generate a private and public key for the server:

     ```shell
     umask 077
     wg genkey | tee /etc/wireguard/privatekey | wg pubkey > /etc/wireguard/publickey
     ```

2. **View the Generated Keys:**
   - You can view the keys using the following commands:

     ```shell
     cat /etc/wireguard/privatekey
     cat /etc/wireguard/publickey
     ```

   - **Important:** Keep the private key secure and do not share it.

## Step 3: Configure WireGuard Interface on OpenWRT

1. **Access LuCI (Web Interface):**
   - Go to the OpenWRT web interface (LuCI) in your browser.
   - Navigate to `Network` > `Interfaces`.

2. **Add a New Interface:**
   - Click on `Add new interface…`.
   - Name the new interface (e.g., `WG0`).
   - Select `WireGuard VPN` as the protocol.

3. **Configure the Interface:**
   - **Private Key:** Enter the private key you generated earlier.
   - **Listen Port:** Set the port you want WireGuard to listen on (e.g., `51820`).
   - **IP Address:** Assign an IP address for the WireGuard interface (e.g., `10.0.0.1/24`).

4. **Save and Apply Changes:**
   - Click `Save` and `Save & Apply`.

## Step 4: Set Up Firewall Rules

1. **Create a Firewall Zone for WireGuard:**
   - Go to `Network` > `Firewall`.
   - Click `Add` to create a new zone.
   - Name the zone (e.g., `wg`).
   - Assign the WireGuard interface (`WG0`) to this new zone.
   - Set the `Input`, `Output`, and `Forward` policies to `ACCEPT`.

2. **Allow Forwarding to LAN:**
   - Under the `Inter-Zone Forwarding` section, check the box to allow forwarding from `wg` to `lan` (and vice versa if needed).

3. **Save and Apply Changes:**
   - Click `Save` and `Save & Apply`.

## Step 5: Add Peer (Client) Configuration

1. **Create Client Keys (on Client Device or OpenWRT):**
   - You can generate client keys on the OpenWRT router or the client device itself:

     ```shell
     umask 077
     wg genkey | tee /etc/wireguard/client_privatekey | wg pubkey > /etc/wireguard/client_publickey
     ```

2. **Add Client (Peer) Configuration in OpenWRT:**
   - In the OpenWRT web interface, go to `Network` > `Interfaces`.
   - Click on the WireGuard interface (`WG0`) and go to the `Peers` tab.
   - Click `Add peer` and fill in the following:
     - **Description:** Name the peer (e.g., `Client1`).
     - **Public Key:** Enter the client’s public key.
     - **Allowed IPs:** Specify the IP address for the client (e.g., `10.0.0.2/32`).
     - **Persistent Keepalive:** Set to `25` seconds (optional, helps with NAT traversal).

3. **Save and Apply Changes:**
   - Click `Save` and `Save & Apply`.

## Step 6: Configure Client Device

1. **Generate or Use Existing Client Configuration:**
   - Use the keys generated earlier for the client.
   - Create a client configuration file (`wg0.conf`):

     ```ini
     [Interface]
     PrivateKey = CLIENT_PRIVATE_KEY
     Address = 10.0.0.X/32
     DNS = OpenWRT_ROUTER_GATEWAY_IP

     [Peer]
     PublicKey = SERVER_PUBLIC_KEY
     Endpoint = SERVER_IP:51820
     AllowedIPs = 0.0.0.0/0
     ```

2. **Load the Configuration on the Client Device:**
   - On the client device, import the configuration file using the WireGuard app or CLI.

3. **Connect the Client:**
   - Activate the WireGuard connection on the client device.

## Step 7: Test the VPN Connection

1. **Check the Connection:**
   - Ensure the client is connected to the WireGuard VPN.
   - On the OpenWRT router, verify the connection by checking the status:

     ```shell
     wg show
     ```

2. **Test Internet Access:**
   - On the client device, check if you can access the internet through the VPN.
   - Use `ipconfig` (Windows) or `ifconfig` (macOS/Linux) to confirm the IP address is from the VPN range.

## Step 8: Persistent Configuration

1. **Make Configuration Persistent:**
   - Ensure that the WireGuard interface starts automatically with the router by adding the following line to `/etc/rc.local`:

     ```shell
     /etc/init.d/network restart
     ```

2. **Reboot the Router:**
   - Reboot the router to ensure all settings are applied and persistent.
