source /etc/kolla/admin-openrc.sh
openstack service list
openstack compute service list
openstack network agent list
openstack token issue

docker exec -it rabbitmq rabbitmqctl cluster_status
docker exec -it openvswitch_vswitchd ovs-vsctl show


openstack network list
openstack image list
openstack hypervisor list
openstack compute service list
openstack server show test-vm-1 -c status -c OS-EXT-STS:vm_state -c OS-EXT-STS:power_state

