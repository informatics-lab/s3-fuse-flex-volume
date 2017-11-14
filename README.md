# NFS FlexVolume

An out of tree, FlexVolume based NFS volume provider for kubernetes.

## What does it do?

It mounts NFS volumes onto pods in a more opinionated way than the built
in NFS volume. It provides the following guarantees / limitations:

1. No NFS export shall be mounted more than once per node
2. Integrated `subPath` support

This makes it perform a lot better for cases where a lot of pods are using
subPath to mount to various directories in the same NFS export. The in-built volume provider in this case does 1 full NFS mount per pod, which is expensive. Instead, we get away with 1 symlink per pod and 1 NFS mount per node.

## When should I use it?

You should use this when your pods to nfs exports ratio is in the hundreds or thousands. 

## How does it work?

1. If this share has not been mounted before with these options, we mount it. We pick a special path to do this to ease housekeeping & testing if this is mounted in the future:

   ```
   <prefix>/host/<host>/path/<escaped-path>/opts/<sorted-opts>
   ```
   
2. We symlink the subPath required to the mount directory specified by kubelet. This makes the contents available to the pod! 

   We can also create subPath if it doesn't exist (if `createIfNecessary` is true).

3. To unmount, we just unlink the symlink and are done.

4. We should keep track of the number of symlinks pointing to the nfs mount, and unmount it when this gets to 0. This is a TODO item still.
