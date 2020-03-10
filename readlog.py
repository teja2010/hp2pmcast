#!/usr/bin/env python
# coding: utf-8

# In[34]:


(/read, hostlist)
(/hostlist, =, [])

hnlist = []
seqthnmap = dict()

for hostname in hostlist:
    file = open(hostname+".log","w")
    
    for line in file.readlines():
        data = line.strip().split(";")
        seq = int(data[0])
        
        if seq not in seqlist:
            seqlist.append(seq)
            seqthnmap[seq] = dict()
            
        tandhn = data[1:]
        
        for thn in tandhn:
            thnmap = dict()
            t, hn = thn.split(",")
            
            if hn not in hnlist:
                t = int(t)
                hostlist.remove(hn)
                hnlist.append(hn)
                thnmap[hn] = t
        seqthnmap[seq].update(thnmap)
    file.close()
    
file = open('logmap.txt', 'w')
for line in seqthnmap:
    file.write(str(line)+': '+str(seqthnmap[line])+'\n')
file.close

