public class SWrapper implements SPtr{
	Struct_int_v data;
	public Struct_int_v value(){return data;}
	public void set(Struct_int_v x){data.set(x);}
	public SWrapper(Struct_int_v data){this.data= data;}	
	public int Square(){return S.Square(data);}
	public boolean equals(Object o){return o instanceof SPtr && ((SPtr)o).value() == data;} // for hash maps
}
