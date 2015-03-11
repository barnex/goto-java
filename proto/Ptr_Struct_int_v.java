class Ptr_Struct_int_v{

	Struct_int_v value;

	public Ptr_Struct_int_v(Struct_int_v value){
		this.value = value;
	}

	public void set(Struct_int_v value){
		this.value.set(value);
	}

	public boolean equals(Object o){
		return o instanceof Ptr_Struct_int_v && ((Ptr_Struct_int_v)o).value.equals(this.value);
	}
}

