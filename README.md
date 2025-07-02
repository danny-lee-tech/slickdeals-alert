# slickdeals-alert
Send email alerts based on certain slickdeals forum criteria

## Current Capabilities
* Notifies if a post has a rank higher than the config-specified value on the first page of the Hot Deals forum for 0+ rank or higher posts
* Ignores the first rules post
* Either notifies via email or PushBullet channel

## How to Use
./slickdeals-alert \<YML config file location\>

## Output Example
New Slickdeals Alert

Id: thread_title_18412342  
Title: [Bundle] AMD Ryzen 7800X3D CPU + ASUS ROG STRIX B650E-I GAMING WIFI MB + Thermal Paste + 1TB Kingston NV3 NVMe SSD $525 + Free S/H  
Rank: 8  
Category: Computers  
Replies: 1  
View Count: 0  
URL: https://www.slickdeals.net//f/18412342-bundle-amd-ryzen-7800x3d-cpu-asus-rog-strix-b650e-i-gaming-wifi-mb-thermal-paste-1tb-kingston-nv3-nvme-ssd-525-free-s-h

## Future Enhancements/Ideas
* Add other filters other than rank, such as category, reply count or view count
* Add the ability to run multiple filter rules (i.e. one filter to find ranks of 8 or more and one filter to find reply count higher than 10) and combine results in a single notification/email
* Smarter Detection Filters - the idea is that a view count of 700 might not be much around noon, but it is considered a lot at 2AM
