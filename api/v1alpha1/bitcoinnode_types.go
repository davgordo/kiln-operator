/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RPCServer struct {
	CertSecret string `json:"certSecret,omitempty"`
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`
}

// BitcoinNodeSpec defines the desired state of BitcoinNode
type BitcoinNodeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	RPCServer     RPCServer `json:"rpcServer,omitempty"`
	MiningAddress string    `json:"miningAddress,omitempty"`
}

// BitcoinNodeStatus defines the observed state of BitcoinNode
type BitcoinNodeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	BlockCount int64 `json:"blockCount"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BitcoinNode is the Schema for the bitcoinnodes API
type BitcoinNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BitcoinNodeSpec   `json:"spec,omitempty"`
	Status BitcoinNodeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BitcoinNodeList contains a list of BitcoinNode
type BitcoinNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BitcoinNode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BitcoinNode{}, &BitcoinNodeList{})
}
