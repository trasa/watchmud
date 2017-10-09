package mobile

// The Mobile standing in front of you is an Instance of
// its definition. Mobiles of definition 'lizard' are all
// immune to poison, but this instance of 'lizard' is wearing
// a magic hat and has a sword in it's hand. (scary lizard)
type Instance struct {
	InstanceId string
	Definition *Definition
}

func NewInstance(instanceId string, defn *Definition) *Instance {
	return &Instance{
		InstanceId: instanceId,
		Definition: defn,
	}
}

func (i *Instance) Id() string {
	return i.InstanceId
}
