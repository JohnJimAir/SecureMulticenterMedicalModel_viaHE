# Top-Lattigo: a higher-level encapsulation of Lattigo, focusing on the operation related to CKKS.

The CKKS scheme supports the pack technique, which means that one single ciphertext contains multiple data. And usually, in real situation, people have to deal with multiple samples, instead of just one sample.
A sample normally has multiple attributes, it is more convenient to pack one attribute of all samples into one ciphertext, instead of pack all attribute of one sample into one ciphertext.
This module provides some functions to help you dealing with multiple samples easily, when using the pack method said above.